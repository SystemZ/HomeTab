package main

import (
	"github.com/spf13/cobra"
	"github.com/systemz/hometab/internal/model"
	"log"
)

func init() {
	rootCmd.AddCommand(convertCmd)
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert data from zfire",
	Run:   convertExec,
}

func convertExec(cmd *cobra.Command, args []string) {
	log.Println("Connecting to DB...")
	model.InitMysql()
	log.Println("Connected to DB")
	model.ImportZfire("./import")
}
