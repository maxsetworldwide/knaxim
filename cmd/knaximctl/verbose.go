package main

import (
	"flag"
	"fmt"
)

var verboseFlag = flag.Bool("v", false, "verbose flag")

func vPrintf(msg string, args ...interface{}) (int, error) {
	if *verboseFlag {
		return fmt.Printf(msg, args...)
	}
	return 0, nil
}
