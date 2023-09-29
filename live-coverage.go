package main

import (
	"log"
	"os"
)

func main() {

	args := os.Args
	if len(args) < 3 {
		log.Fatal("Usage: live-coverage <path-to-prod-coverage> <path-to-test-coverage>")
	}
	prodCoverageFile := args[1]
	testCoverageFile := args[2]

	err := RecountCoverage(prodCoverageFile, testCoverageFile, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
