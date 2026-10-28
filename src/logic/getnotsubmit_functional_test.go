//go:build functional

package logic

import (
	"manabacli/src/manabaclient"
	"net/http/cookiejar"
	"os"
	"testing"
)

func TestGetNotSubmittedAssignments(t *testing.T) {
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

	assignments, err := GetNotSubmittedAssignments(jar)
	if err != nil {
		t.Fatalf("GetNotSubmittedAssignments failed: %v", err)
	}

	for _, a := range assignments {
		t.Logf("[%s] %s -> %s", a.Type, a.Title, a.Url)
	}
}
