package cmd

import (
	"fmt"
	"manabacli/src/config"
	"manabacli/src/logic"
	"manabacli/src/manabaclient"
	"net/http/cookiejar"

	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit <URL> <ファイル名>...",
	Short: "既存のアップロード済みファイルを削除し, 指定したファイルをアップロードして提出します",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		filePaths := args[1:]

		username, password, err := config.Load()
		if err != nil {
			return err
		}

		if err := submitFiles(manabaclient.Client{}, username, password, url, filePaths); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}

type submitter interface {
	manaba
	DeleteAllFiles(jar *cookiejar.Jar, url string) error
}

func submitFiles(m submitter, username, password, url string, filePaths []string) error {
	if err := validate(filePaths); err != nil {
		return err
	}

	jar, err := login(m, username, password)
	if err != nil {
		return err
	}

	// 指定したファイルだけを提出するため, アップロード済みのファイルを先に削除する
	if err := m.DeleteAllFiles(jar, url); err != nil {
		return fmt.Errorf("既存ファイルの削除に失敗しました: %v", err)
	}

	if err := uploadAll(m, jar, url, filePaths); err != nil {
		return err
	}

	if err := logic.SubmitReport(jar, url); err != nil {
		return err
	}
	fmt.Printf("提出に成功しました: %s\n", url)
	return nil
}
