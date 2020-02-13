package cmd

import (
	"github.com/carlogit/phash"
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

	// get files from SQLite
	for page := 0; ; page++ {
		log.Printf("Page %v", page)
		_, imgs := model.List(db, int64(page))
		//_, imgs := model.FileTagSearchByName(db, int64(page), "belly")

		// limit results
		//if page > 1 {
		//	log.Println("Max page reached, ending work")
		//	break
		//}

		// finish if no files left
		if len(imgs) < 1 {
			log.Println("No results left, ending work")
			break
		}

		// scan all entries in DB
		for _, img := range imgs {
			//log.Printf("%v", img.Name)

			// add file
			model.GraphAddFile(dg, model.GraphFile{
				Name:   img.Name,
				Path:   img.Path,
				Sha256: img.Sha256,
				Size:   img.Size,
				Phash:  img.Phash,
				// created?
			})

			// add tags to file
			_, tags := model.TagList(db, img.Fid)
			for _, tag := range tags {
				//log.Printf("tag: %+v", tag)
				// FIXME use UID or something
				model.GraphSetTag(dg, tag.Name, img.Name)
			}

			// set mime for file
			// FIXME use UID or something
			model.GraphSetMime(dg, img.Mime, img.Name)

			// if this is not image, skip distance calc
			if len(img.Phash) < 1 {
				continue
			}
			// set distance
			files := model.GraphSearchPhash(dg)
			for _, file := range files {
				if file.Sha256 == img.Sha256 {
					// don't calc yourself, fool
					continue
				}
				//log.Printf("%+v", file.Uid)
				//log.Printf("%+v", file.Name)
				//log.Printf("%+v", file.Phash)
				distance := phash.GetDistance(img.Phash, file.Phash)
				//log.Printf("dist: %v", distance)
				model.GraphSetDistance(dg, img.Name, file.Name, distance)
			}
		}

	}
}
