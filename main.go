package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	_ "github.com/mattn/go-sqlite3"
)

const (
	redColor    = "\033[31m"
	greenColor  = "\033[32m"
	resetColor  = "\033[0m"
	yellowColor = "\033[33m"
)

type Options struct {
	color       string
	sortBy      string
	sqlFilter   string
	minDiff     int64
	noNew       bool
	prettyPrint bool
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

func (d DiffEntry) String() string {
	return fmt.Sprintf("DiffEntry{filename: %q, diff: %d, after: %d, before: %d}",
		d.filename, d.diff, d.after, d.before)
}

// Function to format numbers with commas
func formatNumber(n int64) string {
	return strconv.FormatInt(int64(n), 10)
}

func setupDatabase(entries []DiffEntry) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Create table
	_, err = db.Exec(`
		CREATE TABLE entries (
			filename TEXT,
			diff INTEGER,
			after INTEGER,
			before INTEGER
		)
	`)
	if err != nil {
		return nil, err
	}

	// Insert data
	stmt, err := db.Prepare("INSERT INTO entries VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for _, entry := range entries {
		_, err = stmt.Exec(entry.filename, entry.diff, entry.after, entry.before)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func main() {
	opts := Options{
		color:  "auto",
		sortBy: "diff",
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <before_file> <after_file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  --color=MODE      color output (MODE: auto, always, never)\n")
		fmt.Fprintf(os.Stderr, "  --sort=TYPE       sort output (TYPE: filename, diff)\n")
		fmt.Fprintf(os.Stderr, "  --min=VALUE       minimum diff value to show (default: math.MinInt64)\n")
		fmt.Fprintf(os.Stderr, "  --no-new          do not include new entries\n")
		fmt.Fprintf(os.Stderr, "  --pretty-print    pretty print numbers\n")
		fmt.Fprintf(os.Stderr, "  --sql=FILTER      SQL WHERE clause for filtering (e.g. 'diff >= 10000 AND filename LIKE '%%session%%')\n")
	}

	colorOpt := flag.String("color", "auto", "")
	sortOpt := flag.String("sort", "diff", "")
	minDiff := flag.Int64("min", math.MinInt64, "minimum diff value to show")
	noNew := flag.Bool("no-new", false, "Do not include new (default: false)")
	prettyPrint := flag.Bool("pretty-print", false, "")
	sqlFilter := flag.String("sql", "", "SQL WHERE clause for filtering (e.g. 'diff > 1000')")

	flag.Parse()
	opts.color = *colorOpt
	opts.sortBy = *sortOpt
	opts.minDiff = *minDiff
	opts.noNew = *noNew
	opts.prettyPrint = *prettyPrint
	opts.sqlFilter = *sqlFilter

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

		if opts.noNew && beforeUsage == 0 {
			continue
		}

		if diff >= opts.minDiff {
			entries = append(entries, DiffEntry{filename, diff, afterUsage, beforeUsage})
		}

	}

	// Add new files that only exist in after
	for _, filename := range afterFiles {
		afterUsage := after[filename]
		if _, exists := before[filename]; !exists {
			diff := afterUsage // beforeUsage is 0
			total += diff

			if !opts.noNew {
				if diff >= opts.minDiff {
					entries = append(entries, DiffEntry{filename, afterUsage, afterUsage, 0})
				}
			}
		}
	}

	// After creating entries slice and before printing:
	db, err := setupDatabase(entries)
	if err != nil {
		fmt.Printf("Error setting up database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Replace existing filtering logic with SQL query
	query := "SELECT filename, diff, after, before FROM entries"
	if opts.sqlFilter != "" {
		query += " WHERE " + opts.sqlFilter
	}
	switch opts.sortBy {
	case "diff":
		query += " ORDER BY diff DESC"
	case "filename":
		query += " ORDER BY filename"
	}

	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	useColor := shouldUseColor(opts.color)
	colorEnd, colorNew := "", ""
	if useColor {
		colorEnd, colorNew = resetColor, yellowColor
	}

	for rows.Next() {
		var entry DiffEntry
		err := rows.Scan(&entry.filename, &entry.diff, &entry.after, &entry.before)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		var beforeStr, afterStr, diffStr string
		if opts.prettyPrint {
			beforeStr = humanize.Comma(entry.before)
			afterStr = humanize.Comma(entry.after)
			diffStr = humanize.Comma(entry.diff)
		} else {
			beforeStr = strconv.FormatInt(entry.before, 10)
			afterStr = strconv.FormatInt(entry.after, 10)
			diffStr = strconv.FormatInt(entry.diff, 10)
		}

		colorStart := ""
		if useColor {
			switch {
			case entry.diff > 0:
				colorStart = redColor
			case entry.diff < 0:
				colorStart = greenColor
			default:
				colorStart = resetColor
			}
		}

		status := ""
		switch {
		case entry.before == entry.after:
			status = "(unchanged)"
		case entry.before == 0:
			status = fmt.Sprintf("%s(new)%s", colorNew, colorEnd)
		case entry.after == 0:
			status = "(removed)"
		}

		fmt.Printf("%-50s | %s%12s (=%12s -%12s) %s%s\n",
			entry.filename,
			colorStart,
			diffStr,
			afterStr,
			beforeStr,
			status,
			colorEnd)
	}
}
