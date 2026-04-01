package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

type NodeResponse struct {
	Hostname        string    `json:"hostname"`
	Status          string    `json:"status"`
	AssignedSSHPort int       `json:"assigned_ssh_port"`
	LastSeenAt      time.Time `json:"last_seen_at"`
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered HIL Benches",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := http.Get(fmt.Sprintf("%s/api/v1/nodes", ServerURL))
		if err != nil {
			fmt.Printf("Error querying central server: %v\n", err)
			os.Exit(1)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			fmt.Printf("Central server returned non-200 status: %d\n", res.StatusCode)
			os.Exit(1)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response API body: %v\n", err)
			os.Exit(1)
		}

		var nodes []NodeResponse
		if err := json.Unmarshal(body, &nodes); err != nil {
			fmt.Printf("Error parsing API JSON: %v\n", err)
			os.Exit(1)
		}

		if len(nodes) == 0 {
			fmt.Println("No active HIL benches found.")
			return
		}

		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		tbl := table.New("Hostname", "Status", "Routing Port", "Last Seen")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, node := range nodes {
			tbl.AddRow(node.Hostname, node.Status, node.AssignedSSHPort, node.LastSeenAt.Format(time.RFC822))
		}

		fmt.Println("\n📡 Active HIL Benches:")
		tbl.Print()
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
