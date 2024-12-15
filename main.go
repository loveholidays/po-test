package main

import (
	"log"
	"os"

	"github.com/sahma19/po-test/pkg/tests"
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
