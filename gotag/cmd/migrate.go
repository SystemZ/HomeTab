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
	Short: "Migrate to new DB",
	Run:   migrateExec,
}

func migrateExec(cmd *cobra.Command, args []string) {
	// get all data and assign them to user ID 1
	model.InitMysql()
	var filesInDb []model.File
	model.DB.Find(&filesInDb)
	log.Printf("Files detected: %v", len(filesInDb))

	// assign all files to user ID 1
	log.Println("Assigning file ownership")
	for _, file := range filesInDb {
		fileUser := model.FileUser{
			FileId:    file.Id,
			UserId:    1,
			CreatedAt: file.CreatedAt,
			UpdatedAt: file.UpdatedAt,
			DeletedAt: nil,
		}
		model.DB.Create(&fileUser)
	}

	// set all tag <-> files connections to user ID 1
	log.Println("Updating file_tags...")
	sql := `UPDATE file_tags SET user_id = 1 WHERE file_tags.user_id IS NULL`
	model.DB.Exec(sql)

	log.Println("All done!")
}
