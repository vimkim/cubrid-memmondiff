package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
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
	rawQuery    string
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

	rootCmd := &cobra.Command{
		Use:   "memmondiff [flags] <before_file> <after_file>",
		Short: "Compare memory snapshots",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
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
			// Rest of your code
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

			var query string
			if opts.rawQuery != "" {
				printRawQuery(db, opts)
				return

			} else {

				// Replace existing filtering logic with SQL query
				query = "SELECT filename, diff, after, before FROM entries"
				if opts.sqlFilter != "" {
					query += " WHERE " + opts.sqlFilter
				}
				switch opts.sortBy {
				case "diff":
					query += " ORDER BY diff DESC"
				case "filename":
					query += " ORDER BY filename"
				}

				fmt.Printf("Query: %s\n", query)

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

				totalStr := strconv.FormatInt(total, 10)

				if opts.prettyPrint {
					totalStr = humanize.Comma(total)
				}

				totalStr = colorize(totalStr, total, opts)
				fmt.Printf("\n# Total Diff: %s\n", totalStr)
			}
		},
	}

	flags := rootCmd.Flags()
	flags.StringVar(&opts.color, "color", "auto", "color output (MODE: auto, always, never)")
	flags.StringVar(&opts.sortBy, "sort", "diff", "(deprecated) sort output (TYPE: filename, diff)")
	flags.Int64Var(&opts.minDiff, "min", math.MinInt64, "(deprecated) minimum diff value to show")
	flags.BoolVar(&opts.noNew, "no-new", false, "(deprecated) do not include new entries")
	flags.BoolVar(&opts.prettyPrint, "pretty-print", false, "pretty print numbers")
	flags.StringVar(&opts.sqlFilter, "sql", "", "SQL WHERE clause for filtering")
	flags.StringVar(&opts.rawQuery, "raw-query", "", "(experimental) Type Raw SQL Query for full control")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func colorize(text string, number int64, opts Options) string {
	useColor := shouldUseColor(opts.color)
	if !useColor {
		return text
	}
	colorEnd := resetColor
	colorStart := ""
	switch {
	case number > 0:
		colorStart = redColor
	case number < 0:
		colorStart = greenColor
	default:
		colorStart = resetColor
	}
	return fmt.Sprintf("%s%s%s", colorStart, text, colorEnd)
}

func printRawQuery(db *sql.DB, opts Options) error {
	rows, err := db.Query(opts.rawQuery)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("getting columns: %w", err)
	}

	// Create result holders
	values := make([]interface{}, len(cols))
	row := make([]interface{}, len(cols))
	for i := range values {
		row[i] = &values[i]
	}

	// Print header
	fmt.Println(strings.Join(cols, "|"))
	fmt.Println(strings.Repeat("-", len(strings.Join(cols, "|"))))

	// Print rows
	for rows.Next() {
		if err := rows.Scan(row...); err != nil {
			return fmt.Errorf("scanning row: %w", err)
		}

		result := make([]string, len(cols))
		for i, val := range values {
			if val == nil {
				result[i] = "NULL"
				continue
			}
			result[i] = fmt.Sprint(val)
		}
		fmt.Println(strings.Join(result, "|"))
	}

	return rows.Err()
}
