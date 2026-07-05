/*
	設定ファイルを読み込むため.
	設定ファイルはusernameとpasswordを保持し, manabaへのログインに使う.

	設定ファイルの形式は次の通り.
	```
	username=yourstudentID
	password=yourpassword
	```
*/

package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// 設定ファイルを読み込んで, usernameとpasswordを取得するため.
func LoadConfig(path string) (username, password string, err error) {
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
