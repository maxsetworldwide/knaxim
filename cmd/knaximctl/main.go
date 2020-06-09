package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println(helpstr)
		return
	}
	switch flag.Arg(0) {
	case "help":
		if flag.NArg() <= 1 {

		}
	default:
		fmt.Println(helpstr)
		return
	}
}
