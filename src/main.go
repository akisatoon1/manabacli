// Command manabacli は、自作ライブラリ github.com/akisatoon1/manaba を使って
// manaba にファイルをアップロードする CLI です。
//
// 使い方:
//
//	manabacli <URL> <ファイル名>...
//
// 認証情報は ~/.manabacli/config から読み込みます。
package main

import (
	"fmt"
	"manabacli/src/config"
	"net/http/cookiejar"
	"os"

	"github.com/akisatoon1/manaba"
)

// usage は引数を誤ったときに表示する使い方の説明を返します。
func usage() string {
	return `使い方: manabacli <URL> <ファイル名>...

  <URL>         提出先の manaba ページ URL
  <ファイル名>  manaba にアップロードするファイルのパス（複数指定可）

認証情報は ~/.manabacli/config に以下の形式で保存してください:
  username=ユーザー名
  password=パスワード`
}

// fatal はエラーメッセージを標準エラー出力に書き出し、指定の終了コードで終了します。
func fatal(code int, format string, args ...any) {
	fmt.Fprintf(os.Stderr, "エラー: "+format+"\n", args...)
	os.Exit(code)
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, usage())
		os.Exit(2)
	}
	url := args[0]
	filePaths := args[1:]

	username, password, err := config.Load()
	if err != nil {
		fatal(1, "設定の読み込みに失敗しました: %v", err)
	}

	uploadFiles(nil, username, password, url, filePaths)
}

func uploadFiles(jar *cookiejar.Jar, username, password, url string, filePaths []string) error {
	// アップロード対象ファイルの存在確認（途中失敗を避けるため、全ファイルを先に検証する）
	for _, filePath := range filePaths {
		info, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				fatal(1, "ファイルが見つかりません: %s", filePath)
			}
			fatal(1, "ファイルを確認できません: %v", err)
		}
		if info.IsDir() {
			fatal(1, "%s はディレクトリです。ファイルを指定してください", filePath)
		}
	}

	// Cookie ジャーを用意してログイン
	jar, err := cookiejar.New(nil)
	if err != nil {
		fatal(1, "Cookie ジャーの作成に失敗しました: %v", err)
	}
	if err := manaba.Login(jar, username, password); err != nil {
		fatal(1, "ログインに失敗しました: %v", err)
	}

	// ファイルを順番にアップロード
	for _, filePath := range filePaths {
		if err := manaba.UploadFile(jar, url, filePath); err != nil {
			fatal(1, "アップロードに失敗しました (%s): %v", filePath, err)
		}
		fmt.Printf("アップロードに成功しました: %s\n", filePath)
	}
	return nil
}
