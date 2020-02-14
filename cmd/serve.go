package cmd

import (
	"gitlab.com/systemz/gotag/model2"
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
	// old interface and SQLite
	db := model.DbInit()
	allFiles := model.CountAllFiles(db)
	log.Printf("All files in DB: %d \n", allFiles)
	go web.Server(db)

	// new interface and MySQL
	model2.InitMysql()
	web.StartWebInterface()
}
