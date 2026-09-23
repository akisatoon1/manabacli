package logic

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// レポート一覧ページを取得し, {レポート名 -> URL} のマップを返す.
// URLを知らなくてもアップロードできるようにするため.
func GetReportNameToUrl(jar *cookiejar.Jar, reportListUrl string) (map[string]string, error) {
	base, err := url.Parse(reportListUrl)
	if err != nil {
		return nil, fmt.Errorf("URL(%v) の解析に失敗しました: %v", reportListUrl, err)
	}

	client := &http.Client{Jar: jar}
	resp, err := client.Get(reportListUrl)
	if err != nil {
		return nil, fmt.Errorf("URL(%v) への GETリクエストに失敗しました: %v", reportListUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("URL(%v) への GETリクエストが成功しませんでした: ステータス: %s", reportListUrl, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTML の解析に失敗しました: %v", err)
	}

	reportNameToUrl, err := createReportNameToUrlFromHTML(doc, base)
	return reportNameToUrl, err
}

// htmlを解析して, {レポート名 -> URL} のマップを作る.
func createReportNameToUrlFromHTML(doc *goquery.Document, base *url.URL) (map[string]string, error) {
	reportNameToUrl := make(map[string]string)

	// レポート一覧テーブルの各行にある <h3 class="report-title"><a href="...">レポート名</a></h3> を拾う
	var err error = nil
	doc.Find("table.stdlist h3.report-title a").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		name := strings.TrimSpace(s.Text())
		if name == "" {
			err = fmt.Errorf("レポート名が空です")
			return false
		}
		href, ok := s.Attr("href")
		if !ok {
			err = fmt.Errorf("レポート(%v) に href 属性がありません", name)
			return false
		}
		ref, parseErr := url.Parse(strings.TrimSpace(href))
		if parseErr != nil {
			err = fmt.Errorf("レポート(%v) の href(%v) の解析に失敗しました: %v", name, href, parseErr)
			return false
		}
		// href は相対 URL なので、一覧ページの URL を基準に絶対 URL へ解決する
		reportNameToUrl[name] = base.ResolveReference(ref).String()
		return true
	})
	if err != nil {
		return nil, err
	}

	return reportNameToUrl, nil
}
