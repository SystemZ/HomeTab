package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/gotag/core"
	"gitlab.com/systemz/gotag/model"
)

func init() {
	rootCmd.AddCommand(diskScan)
}

var diskScan = &cobra.Command{
	Use:   "scan",
	Short: "Scan files on disk",
	Run:   diskScanExec,
}

func diskScanExec(cmd *cobra.Command, args []string) {
	// DB stuff
	db := model.DbInit()
	core.ScanNg(db, args[0])
}
