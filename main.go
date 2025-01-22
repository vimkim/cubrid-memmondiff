package main

import (
	"bufio"
	"flag"
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

var (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
	colors = flag.String("color", "auto", "colorize output (always|auto|never)")
)

func shouldColorize() bool {
	switch *colors {
	case "always":
		return true
	case "never":
		return false
	default:
		fileInfo, _ := os.Stdout.Stat()
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}
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

func colorize(text string, diff int, useColor bool) string {
	if !useColor {
		return text
	}
	if diff > 0 {
		return Red + text + Reset
	}
	return Green + text + Reset
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--color=always|auto|never] <file1> <file2>\n", os.Args[0])
		os.Exit(1)
	}

	entries1, err := readMemFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", args[0], err)
		os.Exit(1)
	}

	entries2, err := readMemFile(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", args[1], err)
		os.Exit(1)
	}

	useColor := shouldColorize()

	for key, entry1 := range entries1 {
		if entry2, exists := entries2[key]; exists {
			if entry1.bytes != entry2.bytes {
				diff := entry2.bytes - entry1.bytes
				output := fmt.Sprintf("%s:%d | %d (=%d-%d)",
					entry1.file, entry1.line, diff, entry2.bytes, entry1.bytes)
				fmt.Println(colorize(output, diff, useColor))
			}
		}
	}
}
