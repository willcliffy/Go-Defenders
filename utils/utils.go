package utils

import (
	"bufio"
	"log"
	"os"
)

func LoadFile(filename string) []string {
	readFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
 
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	readFile.Close()

	return lines
}
