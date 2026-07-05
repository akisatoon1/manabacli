package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// 設定ファイル ~/.manabacli/config の絶対パスを返す.
// パスに~が含まれる不安定さを解消するため.
func GetPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ホームディレクトリを取得できません: %w", err)
	}
	return filepath.Join(home, ".manabacli/config"), nil
}
