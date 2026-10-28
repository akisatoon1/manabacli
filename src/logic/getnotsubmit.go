package logic

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// manaba の未提出課題一覧ページの URL.
const notSubmittedListUrl = "https://room.chuo-u.ac.jp/ct/home_library_query"

// 未提出課題一覧の1行分.
type Assignment struct {
	Type  string // レポート, アンケート, 小テスト など
	Title string
	Url   string // 提出先 URL
}

// 未提出課題一覧ページを取得し, 課題のリストをページの並び順で返す.
// 課題名はコースをまたいで重複し得るため, map ではなくスライスで返す.
func GetNotSubmittedAssignments(jar *cookiejar.Jar) ([]Assignment, error) {
	base, err := url.Parse(notSubmittedListUrl)
	if err != nil {
		return nil, fmt.Errorf("URL(%v) の解析に失敗しました: %v", notSubmittedListUrl, err)
	}

	client := &http.Client{Jar: jar}
	resp, err := client.Get(notSubmittedListUrl)
	if err != nil {
		return nil, fmt.Errorf("URL(%v) への GETリクエストに失敗しました: %v", notSubmittedListUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("URL(%v) への GETリクエストが成功しませんでした: ステータス: %s", notSubmittedListUrl, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTML の解析に失敗しました: %v", err)
	}

	return createAssignmentsFromHTML(doc, base)
}

// htmlを解析して, 課題のリストを作る.
func createAssignmentsFromHTML(doc *goquery.Document, base *url.URL) ([]Assignment, error) {
	assignments := []Assignment{}

	// 一覧テーブルの各行は <td>タイプ</td><td><div class="myassignments-title"><a href="...">タイトル</a></div></td>... という形.
	// 見出し行 (class="title") は飛ばす.
	var err error = nil
	doc.Find("table.stdlist tr").Not(".title").EachWithBreak(func(_ int, row *goquery.Selection) bool {
		typ := strings.TrimSpace(row.Find("td").First().Text())
		link := row.Find(".myassignments-title a").First()
		title := strings.TrimSpace(link.Text())
		if typ == "" || title == "" {
			err = fmt.Errorf("課題のタイプまたはタイトルが空です")
			return false
		}
		href, ok := link.Attr("href")
		if !ok {
			err = fmt.Errorf("課題(%v) に href 属性がありません", title)
			return false
		}
		ref, parseErr := url.Parse(strings.TrimSpace(href))
		if parseErr != nil {
			err = fmt.Errorf("課題(%v) の href(%v) の解析に失敗しました: %v", title, href, parseErr)
			return false
		}
		// href は相対 URL なので、一覧ページの URL を基準に絶対 URL へ解決する
		assignments = append(assignments, Assignment{
			Type:  typ,
			Title: title,
			Url:   base.ResolveReference(ref).String(),
		})
		return true
	})
	if err != nil {
		return nil, err
	}

	return assignments, nil
}
