//go:build functional

package logic

import (
	"manabacli/src/manabaclient"
	"net/http/cookiejar"
	"os"
	"testing"
)

func TestGetReportNameToUrl(t *testing.T) {
	username := os.Getenv("MANABA_USERNAME")
	password := os.Getenv("MANABA_PASSWORD")
	if username == "" || password == "" {
		t.Fatal("MANABA_USERNAME or MANABA_PASSWORD environment variable is not set")
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("Failed to create cookie jar: %v", err)
	}

	client := manabaclient.Client{}
	err = client.Login(jar, username, password)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	reportListUrl := "https://room.chuo-u.ac.jp/ct/course_6345152_report"
	reportNameToUrl, err := GetReportNameToUrl(jar, reportListUrl)
	if err != nil {
		t.Fatalf("GetReportNameToUrl failed: %v", err)
	}

	if len(reportNameToUrl) == 0 {
		t.Fatalf("No reports found")
	}
	for name, url := range reportNameToUrl {
		t.Logf("Report: %s -> URL: %s", name, url)
	}
}
