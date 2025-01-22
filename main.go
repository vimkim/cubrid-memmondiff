package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type MemEntry struct {
	file    string
	line    int
	bytes   int
	percent int
}

func parseMemLine(line string) (*MemEntry, error) {
	re := regexp.MustCompile(`(.*):(\d+) \| (\d+) Bytes\( *(\d+)%\)`)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("invalid format: %s", line)
	}

	lineNum, _ := strconv.Atoi(matches[2])
	bytes, _ := strconv.Atoi(matches[3])
	percent, _ := strconv.Atoi(matches[4])

	return &MemEntry{
		file:    matches[1],
		line:    lineNum,
		bytes:   bytes,
		percent: percent,
	}, nil
}

func readMemFile(filename string) (map[string]*MemEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries := make(map[string]*MemEntry)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Bytes") {
			entry, err := parseMemLine(line)
			if err != nil {
				continue
			}
			key := fmt.Sprintf("%s:%d", entry.file, entry.line)
			entries[key] = entry
		}
	}
	return entries, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file1> <file2>\n", os.Args[0])
		os.Exit(1)
	}

	entries1, err := readMemFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	entries2, err := readMemFile(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", os.Args[2], err)
		os.Exit(1)
	}

	for key, entry1 := range entries1 {
		if entry2, exists := entries2[key]; exists {
			if entry1.bytes != entry2.bytes {
				diff := entry2.bytes - entry1.bytes
				fmt.Printf("%s:%d | %d (=%d-%d)\n",
					entry1.file, entry1.line, diff, entry2.bytes, entry1.bytes)
			}
		}
	}
}
