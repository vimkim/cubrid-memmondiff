package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	redColor    = "\033[31m"
	greenColor  = "\033[32m"
	resetColor  = "\033[0m"
	yellowColor = "\033[33m"
)

type Options struct {
	color   string
	sortBy  string
	minDiff int64
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

func parseMemoryFile(filepath string) (map[string]int64, []string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	result := make(map[string]int64)
	order := make([]string, 0)

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
			order = append(order, filename)
		}
	}
	return result, order, scanner.Err()
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
	minDiff := flag.Int64("min", math.MinInt64, "minimum diff value to show")

	flag.Parse()
	opts.color = *colorOpt
	opts.sortBy = *sortOpt
	opts.minDiff = *minDiff

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	beforeFile := args[0]
	afterFile := args[1]

	before, beforeFiles, err := parseMemoryFile(beforeFile)
	if err != nil {
		fmt.Printf("Error reading before file: %v\n", err)
		os.Exit(1)
	}

	after, afterFiles, err := parseMemoryFile(afterFile)
	if err != nil {
		fmt.Printf("Error reading after file: %v\n", err)
		os.Exit(1)
	}

	entries := make([]DiffEntry, 0)

	var total int64

	// Process files that exist in before
	for _, filename := range beforeFiles {
		beforeUsage := before[filename]
		afterUsage := after[filename]
		diff := afterUsage - beforeUsage
		total += diff

		if diff >= *minDiff {
			entries = append(entries, DiffEntry{filename, diff, afterUsage, beforeUsage})
		}

	}

	// Add new files that only exist in after
	for _, filename := range afterFiles {
		afterUsage := after[filename]
		if _, exists := before[filename]; !exists {
			diff := afterUsage // beforeUsage is 0
			total += diff

			if diff >= *minDiff {
				entries = append(entries, DiffEntry{filename, afterUsage, afterUsage, 0})
			}
		}
	}

	// sort by filename first

	switch opts.sortBy {
	case "diff":
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].diff > entries[j].diff
		})
	case "filename":
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].filename < entries[j].filename
		})
	}

	useColor := shouldUseColor(opts.color)

	colorStart, colorEnd, colorNew := "", "", ""
	if useColor {
		colorEnd, colorNew = resetColor, yellowColor
	}

	for _, entry := range entries {
		if useColor {
			if entry.diff > 0 {
				colorStart = redColor
			} else if entry.diff < 0 {
				colorStart = greenColor
			} else {
				colorStart = resetColor
			}
		}
		if entry.before == entry.after {
			fmt.Printf("%s | %s%d (=%d-%d) (unchanged)%s\n", entry.filename, colorStart, entry.diff, entry.after, entry.before, colorEnd)
		} else if entry.before == 0 {
			fmt.Printf("%s | %s%d (=%d-%d) %s(new)%s\n", entry.filename, colorStart, entry.diff, entry.after, entry.before, colorNew, colorEnd)
		} else if entry.after == 0 {
			fmt.Printf("%s | %s%d (=%d-%d) (removed)%s\n", entry.filename, colorStart, entry.diff, entry.after, entry.before, colorEnd)
		} else {
			fmt.Printf("%s | %s%d (=%d-%d) %s\n", entry.filename, colorStart, entry.diff, entry.after, entry.before, colorEnd)
		}
	}

	if useColor {
		if total > 0 {
			colorStart = redColor
		} else if total < 0 {
			colorStart = greenColor
		}
	}
	fmt.Printf("Total diff: %s%d%s\n", colorStart, total, colorEnd)
}
