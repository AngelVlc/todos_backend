//+build !e2e

package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("TESTING", "true")
	os.Exit(m.Run())
}
