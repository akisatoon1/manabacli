//go:build functional

package logic

import (
	"manabacli/src/manabaclient"
	"net/http/cookiejar"
	"os"
	"testing"
)

func TestSubmitReport(t *testing.T) {
	username := os.Getenv("MANABA_USERNAME")
	password := os.Getenv("MANABA_PASSWORD")
	if username == "" || password == "" {
		t.Fatal("MANABA_USERNAME or MANABA_PASSWORD environment variable is not set")
	}

	reportUrl := "https://room.chuo-u.ac.jp/ct/course_6345179_report_6844361"
	if reportUrl == "" {
		t.Skip("reportUrl is not set")
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

	err = SubmitReport(jar, reportUrl)
	if err != nil {
		t.Fatalf("SubmitReport failed: %v", err)
	}
}
