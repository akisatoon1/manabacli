package main

import (
	"fmt"
	"net/http/cookiejar"
	"os"
)

type manaba interface {
	Login(jar *cookiejar.Jar, username string, password string) error
	UploadFile(jar *cookiejar.Jar, url string, filePath string) error
}

func uploadFiles(m manaba, username, password, url string, filePaths []string) error {
	if err := validate(filePaths); err != nil {
		return err
	}

	// Cookie ジャーを用意してログイン
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("Cookie ジャーの作成に失敗しました: %v", err)
	}
	if err := m.Login(jar, username, password); err != nil {
		return fmt.Errorf("ログインに失敗しました: %v", err)
	}

	// ファイルを順番にアップロード
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
