package logic

import (
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func newDocument(t *testing.T, html string) *goquery.Document {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}
	return doc
}

func TestCreateAssignmentsFromHTML(t *testing.T) {
	html, err := os.ReadFile("testdata/home_library_query.html")
	if err != nil {
		t.Fatalf("Failed to read testdata: %v", err)
	}
	base, _ := url.Parse(notSubmittedListUrl)

	got, err := createAssignmentsFromHTML(newDocument(t, string(html)), base)
	if err != nil {
		t.Fatalf("createAssignmentsFromHTML failed: %v", err)
	}

	want := []Assignment{
		{Type: "レポート", Title: "レポート課題１", Url: "https://room.chuo-u.ac.jp/ct/course_1001_report_2001"},
		{Type: "アンケート", Title: "授業アンケート", Url: "https://room.chuo-u.ac.jp/ct/course_1002_survey_2002"},
		{Type: "小テスト", Title: "第3回小テスト", Url: "https://room.chuo-u.ac.jp/ct/course_1003_query_2003"},
		{Type: "ドリル", Title: "第3回ドリル", Url: "https://room.chuo-u.ac.jp/ct/course_1003_drill_2004"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCreateAssignmentsFromHTML_Empty(t *testing.T) {
	html := `<table class="stdlist"><tr class=title><th>タイプ</th><th>タイトル</th></tr></table>`
	base, _ := url.Parse(notSubmittedListUrl)

	got, err := createAssignmentsFromHTML(newDocument(t, html), base)
	if err != nil {
		t.Fatalf("createAssignmentsFromHTML failed: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %v, want empty", got)
	}
}

func TestCreateAssignmentsFromHTML_NoHref(t *testing.T) {
	html := `<table class="stdlist">
<tr class=title><th>タイプ</th><th>タイトル</th></tr>
<tr class="row0"><td>レポート</td><td><div class="myassignments-title"><a>レポート課題１</a></div></td></tr>
</table>`
	base, _ := url.Parse(notSubmittedListUrl)

	_, err := createAssignmentsFromHTML(newDocument(t, html), base)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
