package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/systemz/tasktab/model"
)

func init() {
	rootCmd.AddCommand(userCreateCmd)
	userCreateCmd.Flags().StringVar(&username, "username", "", "Username for auth")
	userCreateCmd.MarkFlagRequired("username")
	userCreateCmd.Flags().StringVar(&password, "password", "", "Password for auth")
	userCreateCmd.MarkFlagRequired("password")
	userCreateCmd.Flags().StringVar(&email, "email", "", "Email (not used yet)")
	userCreateCmd.MarkFlagRequired("email")
}

var (
	username string
	password string
	email    string
)

var userCreateCmd = &cobra.Command{
	Use:   "user-create",
	Short: "Create new user account",
	Run:   userCreateCmdExec,
}

func userCreateCmdExec(cmd *cobra.Command, args []string) {
	model.InitMysql()
	model.CreateUser(username, email, password)
}
