package cmd

import (
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
	model.InitMysql()
	model.InitRedis()
	web.StartWebInterface()
}
