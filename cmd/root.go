package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "docker-cleaner",
	Short: "A CLI tool to clean Docker resources",
	Long: `Docker Cleaner is a CLI tool that helps you clean up Docker resources
including containers, images, volumes, and networks to free up disk space.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Docker Cleaner - Use --help for available commands")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}
