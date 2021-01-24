package main

import (
	"github.com/spf13/cobra"
	"github.com/systemz/hometab/internal/model"
	"github.com/systemz/hometab/internal/server"
	"github.com/systemz/hometab/internal/service/cron"
	"log"
	"time"
)

func init() {
	rootCmd.AddCommand(wwwCmd)
}

var wwwCmd = &cobra.Command{
	Use:   "www",
	Short: "HTTP server",
	Long:  `Serves web interface requests`,
	Run:   wwwExec,
}

func wwwExec(cmd *cobra.Command, args []string) {
	log.Println("Wild TaskTab appears!")
	model.InitMysql()
	model.InitRedis()
	// simple background cron task
	go func() {
		for true {
			cron.ScanRecurring(model.DB)
			time.Sleep(time.Second * 30)
		}
	}()

	server.StartWebInterface()
	log.Println("Dying...")
}
