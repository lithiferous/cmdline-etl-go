package cmd

import (
	"github.com/lithiferous/cmd-etl/api"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Upload csv files in the given path to Postgres.",
	Long: `Iterate over a given list of filenames, extract snapshot
	        date from the name of the file and upload them to a database.`,
	Run: func(cmd *cobra.Command, args []string) {
		api.UploadFilesToDB(cmd.Root().Context(), args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
