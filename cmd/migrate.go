package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/gotag/model"
	"log"
)

func init() {
	rootCmd.AddCommand(migrate)
}

var migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate to Graph DB",
	Run:   migrateExec,
}

func migrateExec(cmd *cobra.Command, args []string) {
	// SQLite stuff
	db := model.DbInit()
	allFiles := model.CountAllFiles(db)
	log.Printf("all files: %v", allFiles)

	// dgraph stuff
	dg := model.GraphInit()

	// get few files as sample
	_, imgs := model.FileTagSearchByName(db, 1, "tag")
	for _, img := range imgs {
		log.Printf("%v", img.Name)

		// add file
		model.GraphAddFile(dg, model.GraphFile{
			Name:   img.Name,
			Sha256: img.Sha256,
		})

		// add tags to file
		_, tags := model.TagList(db, img.Fid)
		for _, tag := range tags {
			log.Printf("tag: %+v", tag)
			model.GraphSetTag(dg, tag.Name, img.Name)
		}

	}
}
