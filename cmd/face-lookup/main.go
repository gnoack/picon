package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gnoack/picon"
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatalf("wrong number of arguments, needs one email address")
	}

	f, ok := picon.Lookup(flag.Args()[0])
	if ok {
		fmt.Println(f)
	} else {
		// Not found
		os.Exit(1)
	}
}
