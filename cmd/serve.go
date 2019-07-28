package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/systemz/gotag/model"
	"gitlab.com/systemz/gotag/web"
)

func init() {
	rootCmd.AddCommand(httpServe)
}

var httpServe = &cobra.Command{
	Use:   "serve",
	Short: "Run HTTP server",
	Run:   serveExec,
}

func serveExec(cmd *cobra.Command, args []string) {
	// DB stuff
	db := model.DbInit()
	allFiles := model.CountAllFiles(db)
	log.Printf("All files in DB: %d \n", allFiles)
	web.Server(db)
}
