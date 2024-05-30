/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/lithiferous/cmd-etl/api"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A list of relevant information about imported snapshots.",
	Long: `Search snapshot files in source directory, compare it with
				 current database states, i.e.:
         - Which snapshots were imported?
         - When were the snapshots imported?
         - Is this snapshot file already imported the database?`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = api.FindFilesToUpload(cmd.Root().Context())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
