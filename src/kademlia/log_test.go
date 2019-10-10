package main

import (
	"testing"
)

func TestLog(t *testing.T) {
	Log("test")
	EnableLogs()
	Log("test")
	DisableLogs()
}
