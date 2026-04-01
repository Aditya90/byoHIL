package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [hostname]",
	Short: "Establish an interactive SSH session to a HIL Bench",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetHost := args[0]

		fmt.Printf("🔍 Locating bench '%s' on %s...\n", targetHost, ServerURL)

		// 1. Fetch Node routing port
		res, err := http.Get(fmt.Sprintf("%s/api/v1/nodes", ServerURL))
		if err != nil {
			fmt.Printf("❌ Failed to contact Central Server: %v\n", err)
			os.Exit(1)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("❌ Failed to read server response: %v\n", err)
			os.Exit(1)
		}

		var nodes []NodeResponse
		if err := json.Unmarshal(body, &nodes); err != nil {
			fmt.Printf("❌ Failed to parse node response json: %v\n", err)
			os.Exit(1)
		}

		var targetNode *NodeResponse
		for _, node := range nodes {
			if node.Hostname == targetHost {
				targetNode = &node
				break
			}
		}

		if targetNode == nil {
			fmt.Printf("❌ Node '%s' not found. Ensure it is powered on and registered.\n", targetHost)
			os.Exit(1)
		}

		// 2. Perform Post to AccessLog
		devUser := os.Getenv("USER")
		if devUser == "" {
			devUser = "unknown_cli_user"
		}

		logPayload := map[string]string{
			"username": devUser,
			"action":   "cli_connect",
		}
		logBody, _ := json.Marshal(logPayload)

		logURL := fmt.Sprintf("%s/api/v1/nodes/%s/access_log", ServerURL, targetHost)
		logReq, _ := http.NewRequest("POST", logURL, bytes.NewBuffer(logBody))
		logReq.Header.Set("Content-Type", "application/json")
		
		client := &http.Client{Timeout: 5 * time.Second}
		_, err = client.Do(logReq)
		if err != nil {
			fmt.Printf("⚠️ Warning: Failed to write access log to central server (%v)\n", err)
		}

		// 3. Spawning os/exec for ProxyJump reverse SSH
		fmt.Printf("🚇 Establishing Secure Reverse Tunnel to %s (Port %d)...\n\n", targetNode.Hostname, targetNode.AssignedSSHPort)

		centralSSHTarget := os.Getenv("HIL_SSH_HOST")
		if centralSSHTarget == "" {
			centralSSHTarget = "localhost" // Defaults to localhost for testing
		}

		// Command: ssh -t -J <SSHUser>@centralSSHTarget <SSHUser>@localhost -p <targetPort>
		// Note: The python agent binds its port 22 to the targetPort on the central server loopback.
		sshArgs := []string{
			"-t",  // Force TTY allocation
			"-J", fmt.Sprintf("%s@%s", SSHUser, centralSSHTarget), // ProxyJump through central server
			"-p", fmt.Sprintf("%d", targetNode.AssignedSSHPort),   // Port on central server loopback
			"-o", "StrictHostKeyChecking=no",                      // Ignore bench keys for now
			fmt.Sprintf("%s@localhost", SSHUser),                  // Connecting to the bound loopback endpoint
		}

		sshCmd := exec.Command("ssh", sshArgs...)
		sshCmd.Stdin = os.Stdin
		sshCmd.Stdout = os.Stdout
		sshCmd.Stderr = os.Stderr

		if err := sshCmd.Run(); err != nil {
			fmt.Printf("\n❌ SSH Tunnel collapsed or failed to connect: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("\n✅ Connection Closed.")
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
