package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type CodeBlockInfo struct {
	positions  string
	statements int
	isCovered  bool
}

// RecountCoverage - creates coverage
func RecountCoverage(prodCoverageFile string, testCoverageFile string, w io.Writer) error {
	prodBlocks, err := getBlocks(prodCoverageFile)
	if err != nil {
		return fmt.Errorf("failed to get blocks from prod coverage file: %w", err)
	}
	testBlocks, err := getBlocks(testCoverageFile)
	if err != nil {
		return fmt.Errorf("failed to get blocks from test coverage file: %w", err)
	}

	tested := make(map[string]bool)
	for _, block := range testBlocks {
		tested[block.positions] = block.isCovered
	}

	fmt.Println("mode: set")
	for _, block := range prodBlocks {
		// skip blocks that are not executed in prod
		if !block.isCovered {
			continue
		}

		var isCoveredInt int

		if tested[block.positions] {
			isCoveredInt = 1
		}

		_, err = fmt.Fprintf(w, "%s %d %d\n", block.positions, block.statements, isCoveredInt)
		if err != nil {
			return err
		}
	}

	return nil
}

func getBlocks(fileName string) ([]CodeBlockInfo, error) {
	fileReader, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()
	scanner := bufio.NewScanner(fileReader)
	scanner.Split(bufio.ScanLines)
	var blocks []CodeBlockInfo
	for scanner.Scan() {
		line := scanner.Text()
		if line == "mode: set" {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		positions := parts[0]
		statements, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid line: %s", line)
		}
		isTested := parts[2] == "1"
		blocks = append(blocks, CodeBlockInfo{
			positions:  positions,
			statements: statements,
			isCovered:  isTested,
		})
	}
	return blocks, nil
}
