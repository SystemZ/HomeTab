package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/tasktab/model"
	"log"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data from Zfire",
	Run:   importExec,
}

func importExec(cmd *cobra.Command, args []string) {
	log.Println("Connecting to DB...")
	model.InitMysql()
	log.Println("Connected to DB")
	model.ImportZfire("./import")
}
