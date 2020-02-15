package cmd

import (
	"github.com/spf13/cobra"
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
	//core.ScanNg(db, args[0])
}
