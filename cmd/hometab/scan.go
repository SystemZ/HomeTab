package main

import (
	"github.com/spf13/cobra"
	"github.com/systemz/hometab/internal/model"
	"github.com/systemz/hometab/internal/service/gotagcore"
)

func init() {
	rootCmd.AddCommand(diskScan)
	diskScan.Flags().IntVarP(&scanUserID, "user-id", "u", 0, "User ID")
	diskScan.MarkFlagRequired("user-id")
}

var (
	scanUserID int
)

var diskScan = &cobra.Command{
	Use:   "scan",
	Short: "Add files from disk to DB",
	Args:  cobra.MinimumNArgs(1),
	Run:   diskScanExec,
}

func diskScanExec(cmd *cobra.Command, args []string) {
	model.InitMysql()
	gotagcore.ScanMono(model.DB, args[0], scanUserID)
}
