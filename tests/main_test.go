package tests

import (
	"os"
	"testing"

	"api/config"
)

func TestMain(m *testing.M) {
	config.Config()
	defer config.CloseConfig()
	code := m.Run()
	os.Exit(code)
}
