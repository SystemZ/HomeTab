package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/tasktab/model"
	"gitlab.com/systemz/tasktab/service/cron"
	"gitlab.com/systemz/tasktab/web"
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

	web.StartWebInterface()
	log.Println("Dying...")
}
