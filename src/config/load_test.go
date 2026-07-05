package config_test

import (
	"manabacli/src/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testcase struct {
	name             string
	fileContent      string
	expectedUsername string
	expectedPassword string
}

func TestLoadConfig(t *testing.T) {
	tests := []testcase{
		{
			name:             "successful",
			fileContent:      "username=testun\npassword=testpd\n",
			expectedUsername: "testun",
			expectedPassword: "testpd",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			runTc(tt, tc)
		})
	}
}

func runTc(t *testing.T, tc testcase) {
	// arrange
	dir := t.TempDir()
	path := filepath.Join(dir, "config")
	if err := os.WriteFile(path, []byte(tc.fileContent), 0400); err != nil {
		t.Fatal(err)
	}

	// act
	username, passwd, err := config.LoadConfig(path)

	// assert
	require.NoError(t, err)
	assert.Equal(t, tc.expectedUsername, username)
	assert.Equal(t, tc.expectedPassword, passwd)
}
