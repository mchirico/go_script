package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//a = pkg.App{}

	//a.Initilize()
	code := m.Run()

	os.Exit(code)
}
