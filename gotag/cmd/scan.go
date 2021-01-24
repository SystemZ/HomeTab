package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/gotag/core"
	"gitlab.com/systemz/gotag/model"
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
	db := model.InitMysql()
	core.ScanMono(db, args[0], scanUserID)
}
