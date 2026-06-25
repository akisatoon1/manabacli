// Command manaba は、自作ライブラリ github.com/akisatoon1/manaba を使って
// manaba にファイルをアップロードする CLI です。
//
// 使い方:
//
//	manaba <ファイル名> <URL>
//
// 認証情報は ~/.manabacli/config から読み込みます。
package main

import (
	"bufio"
	"fmt"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strings"

	"github.com/akisatoon1/manaba"
)

const configRelPath = ".manabacli/config"

// usage は引数を誤ったときに表示する使い方の説明を返します。
func usage() string {
	return `使い方: manaba <ファイル名> <URL>

  <ファイル名>  manaba にアップロードするファイルのパス
  <URL>         提出先の manaba ページ URL

認証情報は ~/.manabacli/config に以下の形式で保存してください:
  username=ユーザー名
  password=パスワード`
}

// fatal はエラーメッセージを標準エラー出力に書き出し、指定の終了コードで終了します。
func fatal(code int, format string, args ...any) {
	fmt.Fprintf(os.Stderr, "エラー: "+format+"\n", args...)
	os.Exit(code)
}

// configPath は設定ファイル ~/.manabacli/config の絶対パスを返します。
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ホームディレクトリを取得できません: %w", err)
	}
	return filepath.Join(home, configRelPath), nil
}

// loadConfig は設定ファイルを読み込み、username と password を返します。
// 1行1項目の key=value 形式で、'#' で始まる行と空行は無視します。
// 値に '=' を含められるよう、最初の '=' のみで分割します。
func loadConfig(path string) (username, password string, err error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", fmt.Errorf("設定ファイルが見つかりません: %s\n\n"+
				"以下の内容で作成してください:\n"+
				"  mkdir -p ~/.manabacli\n"+
				"  cat > ~/.manabacli/config <<'EOF'\n"+
				"  username=あなたのユーザー名\n"+
				"  password=あなたのパスワード\n"+
				"  EOF\n"+
				"  chmod 600 ~/.manabacli/config", path)
		}
		return "", "", fmt.Errorf("設定ファイルを開けません: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		switch key {
		case "username":
			username = value
		case "password":
			password = value
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	var missing []string
	if username == "" {
		missing = append(missing, "username")
	}
	if password == "" {
		missing = append(missing, "password")
	}
	if len(missing) > 0 {
		return "", "", fmt.Errorf("設定ファイル %s に %s がありません。\n"+
			"username=... と password=... の両方を記載してください",
			path, strings.Join(missing, " と "))
	}
	return username, password, nil
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, usage())
		os.Exit(2)
	}
	filePath := args[0]
	url := args[1]

	// 設定ファイルの読み込み
	cfgPath, err := configPath()
	if err != nil {
		fatal(1, "%v", err)
	}
	username, password, err := loadConfig(cfgPath)
	if err != nil {
		fatal(1, "%v", err)
	}

	// アップロード対象ファイルの存在確認
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

	// Cookie ジャーを用意してログイン
	jar, err := cookiejar.New(nil)
	if err != nil {
		fatal(1, "Cookie ジャーの作成に失敗しました: %v", err)
	}
	if err := manaba.Login(jar, username, password); err != nil {
		fatal(1, "ログインに失敗しました: %v", err)
	}

	// ファイルをアップロード
	if err := manaba.UploadFile(jar, url, filePath); err != nil {
		fatal(1, "アップロードに失敗しました: %v", err)
	}

	fmt.Printf("アップロードに成功しました: %s\n", filePath)
}
