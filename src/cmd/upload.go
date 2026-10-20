package cmd

import (
	"fmt"
	"manabacli/src/config"
	"manabacli/src/manabaclient"
	"net/http/cookiejar"
	"os"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload <URL> <ファイル名>...",
	Short: "ファイルを manaba にアップロードします",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		filePaths := args[1:]

		username, password, err := config.Load()
		if err != nil {
			return err
		}

		if err := uploadFiles(manabaclient.Client{}, username, password, url, filePaths); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

type manaba interface {
	Login(jar *cookiejar.Jar, username string, password string) error
	UploadFile(jar *cookiejar.Jar, url string, filePath string) error
}

func uploadFiles(m manaba, username, password, url string, filePaths []string) error {
	if err := validate(filePaths); err != nil {
		return err
	}

	jar, err := login(m, username, password)
	if err != nil {
		return err
	}
	return uploadAll(m, jar, url, filePaths)
}

// manabaにファイルをアップロードするために, manabaにログインをする
func login(m manaba, username, password string) (*cookiejar.Jar, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("Cookie ジャーの作成に失敗しました: %v", err)
	}
	if err := m.Login(jar, username, password); err != nil {
		return nil, fmt.Errorf("ログインに失敗しました: %v", err)
	}
	return jar, nil
}

// ファイルを順番にアップロード
func uploadAll(m manaba, jar *cookiejar.Jar, url string, filePaths []string) error {
	for _, filePath := range filePaths {
		if err := m.UploadFile(jar, url, filePath); err != nil {
			return fmt.Errorf("アップロードに失敗しました (%s): %v", filePath, err)
		}
		fmt.Printf("アップロードに成功しました: %s\n", filePath)
	}
	return nil
}

func validate(filePaths []string) error {
	// アップロード対象ファイルの存在確認（途中失敗を避けるため、全ファイルを先に検証する）
	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("ファイルが見つかりません: %s", filePath)
			}
			return fmt.Errorf("ファイルを確認できません: %v", err)
		}
		if info.IsDir() {
			return fmt.Errorf("%s はディレクトリです。ファイルを指定してください", filePath)
		}
	}
	return nil
}
