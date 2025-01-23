package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	redColor   = "\033[31m"
	greenColor = "\033[32m"
	resetColor = "\033[0m"
)

type Options struct {
	color  string
	sortBy string
}

type customFlag struct {
	set   func(string)
	value string
}

func (f *customFlag) Set(value string) error {
	f.set(value)
	return nil
}

func (f *customFlag) String() string {
	return f.value
}

func isTerminal() bool {
	if stdoutInfo, err := os.Stdout.Stat(); err == nil {
		return (stdoutInfo.Mode() & os.ModeCharDevice) != 0
	}
	return false
}

func shouldUseColor(opt string) bool {
	switch opt {
	case "always":
		return true
	case "never":
		return false
	default: // "auto"
		return isTerminal()
	}
}

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

type DiffEntry struct {
	filename string
	diff     int64
	after    int64
	before   int64
}

func main() {
	opts := Options{
		color:  "auto",
		sortBy: "diff",
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <before_file> <after_file>\n\nOptions:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  --color=MODE               color output (MODE: auto, always, never)\n")
		fmt.Fprintf(os.Stderr, "  --sort=TYPE                sort output (TYPE: filename, diff)\n")
	}

	colorOpt := flag.String("color", "auto", "")
	sortOpt := flag.String("sort", "diff", "")

	flag.Parse()
	opts.color = *colorOpt
	opts.sortBy = *sortOpt

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	beforeFile := args[0]
	afterFile := args[1]

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

	var entries []DiffEntry

	for filename, afterUsage := range after {
		if beforeUsage, exists := before[filename]; exists {
			diff := afterUsage - beforeUsage
			if diff != 0 {
				entries = append(entries, DiffEntry{filename, diff, afterUsage, beforeUsage})
			}
		} else {
			entries = append(entries, DiffEntry{filename, afterUsage, afterUsage, 0})
		}
	}

	for filename, beforeUsage := range before {
		if _, exists := after[filename]; !exists {
			entries = append(entries, DiffEntry{filename, -beforeUsage, 0, beforeUsage})
		}
	}

	// sort by filename first
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].filename < entries[j].filename
	})

	switch opts.sortBy {
	case "diff":
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].diff > entries[j].diff // descending
		})
	}

	useColor := shouldUseColor(opts.color)

	for _, entry := range entries {
		color := ""
		if useColor {
			if entry.diff > 0 {
				color = redColor
			} else if entry.diff < 0 {
				color = greenColor
			}
		}

		colorStart, colorEnd := "", ""
		if useColor {
			colorStart, colorEnd = color, resetColor
		}
		if entry.before == 0 {
			fmt.Printf("%s%s | %d (new)%s\n", colorStart, entry.filename, entry.diff, colorEnd)
		} else if entry.after == 0 {
			fmt.Printf("%s%s | %d (removed)%s\n", colorStart, entry.filename, entry.diff, colorEnd)
		} else {
			fmt.Printf("%s%s | %d (=%d-%d)%s\n", colorStart, entry.filename, entry.diff, entry.after, entry.before, colorEnd)
		}
	}
}
