/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/lithiferous/cmd-etl/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/fsnotify/fsnotify"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Watch source csv file directory for new arrivals and upload new files to db",
	Long: `Continously track changes in files for a given directory, if a new file is detected
	       launch upload to a database.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create new watcher.
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal().Err(err)
		}
		defer watcher.Close()

		// Start listening for events.
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					log.Info().Msgf("event: %s", event)
					if event.Has(fsnotify.Create) {
						log.Info().Msgf("modified file: %s", event.Name)
						new_files := api.FindFilesToUpload(cmd.Root().Context())

						if len(new_files) == 0 {
							log.Info().Msg("no new files discovered")
						} else {
							api.UploadFilesToDB(cmd.Context(), new_files)
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Fatal().Err(err)
				}
			}
		}()

		// access context's state to infer environment
		// variable to check for raw files on folder
		state := cmd.Root().Context().Value(api.AppState{}).(*api.State)
		source_dir := (*state).Config.DirSource

		// Add a path.
		err = watcher.Add(source_dir)

		if err != nil {
			log.Fatal().Err(err)
		}

		// Block main goroutine forever.
		<-make(chan struct{})
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
