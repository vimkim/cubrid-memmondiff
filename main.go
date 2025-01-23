package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	redColor   = "\033[31m"
	greenColor = "\033[32m"
	resetColor = "\033[0m"
)

func parseMemoryFile(filepath string) (map[string]int64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]int64)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "|") && !strings.Contains(line, "File Name") {
			parts := strings.Split(line, "|")
			if len(parts) != 2 {
				continue
			}

			filename := strings.TrimSpace(parts[0])
			usageStr := strings.TrimSpace(parts[1])
			usageStr = strings.Split(usageStr, " ")[0]

			usage, err := strconv.ParseInt(usageStr, 10, 64)
			if err != nil {
				continue
			}

			result[filename] = usage
		}
	}

	return result, scanner.Err()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: program <before_file> <after_file>")
		os.Exit(1)
	}

	beforeFile := os.Args[1]
	afterFile := os.Args[2]

	before, err := parseMemoryFile(beforeFile)
	if err != nil {
		fmt.Printf("Error reading before file: %v\n", err)
		os.Exit(1)
	}

	after, err := parseMemoryFile(afterFile)
	if err != nil {
		fmt.Printf("Error reading after file: %v\n", err)
		os.Exit(1)
	}

	diffs := make(map[string]int64)

	for filename, afterUsage := range after {
		if beforeUsage, exists := before[filename]; exists {
			diff := afterUsage - beforeUsage
			if diff != 0 {
				diffs[filename] = diff
			}
		} else {
			diffs[filename] = afterUsage
		}
	}

	for filename, beforeUsage := range before {
		if _, exists := after[filename]; !exists {
			diffs[filename] = -beforeUsage
		}
	}

	for filename, diff := range diffs {
		originalValue := before[filename]
		color := ""
		if diff > 0 {
			color = redColor
		} else if diff < 0 {
			color = greenColor
		}

		if diff > 0 {
			fmt.Printf("%s%s | %d (=%d-%d)%s\n", color, filename, diff, after[filename], originalValue, resetColor)
		} else {
			if _, exists := after[filename]; exists {
				fmt.Printf("%s%s | %d (=%d-%d)%s\n", color, filename, diff, after[filename], originalValue, resetColor)
			} else {
				fmt.Printf("%s%s | %d (removed)%s\n", color, filename, diff, resetColor)
			}
		}
	}
}
