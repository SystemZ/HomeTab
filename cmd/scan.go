package cmd

import (
	"log"

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
	allFiles := model.CountAllFiles(db)
	log.Printf("All files in DB: %d \n", allFiles)
	core.ScanNg(db, args[0])
}
