package logic

import (
	"fmt"
	"net/http/cookiejar"

	"github.com/akisatoon1/manaba"
)

// レポート提出ページの URL に対して提出を行う.
// アップロード済みのファイルをブラウザを開かずに提出できるようにするため.
func SubmitReport(jar *cookiejar.Jar, reportUrl string) error {
	if err := manaba.SubmitReports(jar, reportUrl); err != nil {
		return fmt.Errorf("URL(%v) の提出に失敗しました: %v", reportUrl, err)
	}
	return nil
}
