package controllers

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("TESTING", "true")
	os.Exit(m.Run())
}
