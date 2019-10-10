package main

import (
	"fmt"
)

// EnableLog is a flag is logs are enabled
var EnableLog = false

//Log is used to insert a string into the log.
func Log(info string) {
	if !EnableLog {
		return
	}
	fmt.Println(info)
}

// EnableLogs isa setter for the log flag
func EnableLogs() {
	EnableLog = true
}

// DisableLogs is a setter for disabling the log
func DisableLogs() {
	EnableLog = false
}
