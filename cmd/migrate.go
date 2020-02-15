package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrate)
}

var migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate to new DB",
	Run:   migrateExec,
}

func migrateExec(cmd *cobra.Command, args []string) {
	/*
		mysql := model2.InitMysql()

		// scan all entries in DB
		//imgs := model.ListAll(sqlite)
		for _, img := range imgs {

			log.Printf("%v %v %v", img.Fid, img.Name, img.Sha256)

			//time.Sleep(time.Millisecond * 50)
			//continue

			// upgrade pHash storage
			pHashA := 0
			pHashB := 0
			pHashC := 0
			pHashD := 0
			if len(img.Phash) > 1 {
				pHashA, _ = strconv.Atoi(img.Phash[0:16])
				pHashB, _ = strconv.Atoi(img.Phash[16:32])
				pHashC, _ = strconv.Atoi(img.Phash[32:48])
				pHashD, _ = strconv.Atoi(img.Phash[48:64])
			}

			// save MIME to DB
			mimeId := model2.AddMime(mysql, img.Mime)

			// save file to DB
			file := &model2.File{
				Filename: img.Name,
				FilePath: img.Path,
				SizeB:    img.Size,
				MimeId:   mimeId,
				Sha256:   img.Sha256,
				PhashA:   pHashA,
				PhashB:   pHashB,
				PhashC:   pHashC,
				PhashD:   pHashD,
			}
			mysql.Save(&file)

			// add tags to DB
			found, tags := model.TagList(sqlite, img.Fid)
			if !found {
				// finish work if no tags for this file
				continue
			}
			for _, tag := range tags {
				model2.AddTagToFile(mysql, tag.Name, file.Id)
			}
		}
	*/
}
