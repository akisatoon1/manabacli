package cmd

import (
	"fmt"
	"net/http/cookiejar"

	"manabacli/src/config"
	"manabacli/src/logic"
	"manabacli/src/manabaclient"

	"github.com/spf13/cobra"
)

var getnotsubmitCmd = &cobra.Command{
	Use:   "getnotsubmit",
	Short: "未提出課題の一覧を, タイプ・タイトル・提出先URLの形式で取得します.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, password, err := config.Load()
		if err != nil {
			return err
		}

		jar, err := cookiejar.New(nil)
		if err != nil {
			return err
		}

		client := manabaclient.Client{}
		err = client.Login(jar, username, password)
		if err != nil {
			return err
		}

		assignments, err := logic.GetNotSubmittedAssignments(jar)
		if err != nil {
			return err
		}

		if len(assignments) == 0 {
			fmt.Println("未提出の課題はありません")
			return nil
		}
		for _, a := range assignments {
			fmt.Printf("[%s] %s -> %s\n", a.Type, a.Title, a.Url)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getnotsubmitCmd)
}
