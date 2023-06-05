package main

import (
	"github.com/loveholidays/po-test/pkg/tests"
	"log"
	"os"
)

func main() {

	if len(os.Args[1:]) == 0 {
		log.Fatalf("Usage: po-test test-file1.yaml test-file2.yaml")
	}

	err := tests.RunUnitTests(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
}
