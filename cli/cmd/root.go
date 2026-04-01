package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	ServerURL string
	SSHUser   string
)

var rootCmd = &cobra.Command{
	Use:   "hilcli",
	Short: "Hardware-in-the-Loop CLI Tool",
	Long:  `hilcli is a tool to securely find and connect to physical HIL test node benches via central reverse SSH tunneling.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Let the user override via env var, otherwise default to local development setup
	defaultURL := os.Getenv("HIL_SERVER_URL")
	if defaultURL == "" {
		defaultURL = "http://localhost:8080"
	}
	
	defaultUser := os.Getenv("HIL_SSH_USER")
	if defaultUser == "" {
		// Default to testing against your own mac
		defaultUser = os.Getenv("USER")
	}

	rootCmd.PersistentFlags().StringVarP(&ServerURL, "server", "s", defaultURL, "Central Go Server API URL")
	rootCmd.PersistentFlags().StringVarP(&SSHUser, "user", "u", defaultUser, "SSH Username for the reverse tunnel ProxyJump")
}
