package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "manabacli",
	Short: "manaba に関する処理を行う CLI ツール",
	Long: `使い方: manabacli <URL> <ファイル名>...

  <URL>         提出先の manaba ページ URL
  <ファイル名>  manaba にアップロードするファイルのパス（複数指定可）

認証情報は ~/.manabacli/config に以下の形式で保存してください:
  username=ユーザー名
  password=パスワード`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
