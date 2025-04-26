package tests

import (
	"os"
	"testing"

	"api/config"
)

func TestMain(m *testing.M) {
	config.Config()
	code := m.Run()
	os.Exit(code)
}
