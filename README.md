# memmondiff

Command-line tool for analyzing memory usage changes in CUBRID memmon files using SQL filtering.

## Installation

### Pre-built Binary

```bash
wget https://github.com/vimkim/cubrid-memmondiff/releases/latest/download/memmondiff-linux-amd64
```

### Go Install

```bash
go install github.com/vimkim/cubrid-memmondiff@latest
```

### Build from Source

```bash
git clone https://github.com/vimkim/cubrid-memmondiff
cd cubrid-memmondiff
go build .
```

## Usage

```bash
memmondiff [flags] <before_file> <after_file>
```

### Flags

- `--color`: Output coloring (auto, always, never)
- `--pretty-print`: Format numbers with commas
- `--sql`: Filter results using SQL WHERE clause
- `--raw-query`: Execute arbitrary SQL queries against the diff data

### Basic Examples

```bash
# Compare two memmon files
memmondiff before.txt after.txt

# Pretty print numbers
memmondiff --pretty-print before.txt after.txt

# Filter large memory changes
memmondiff --sql "diff >= 10000" before.txt after.txt
```

### Advanced SQL Filtering

```bash
# Filter by filename pattern
memmondiff --sql "filename LIKE '%heap%'" before.txt after.txt

# Complex conditions
memmondiff --sql "diff >= 5000 AND filename NOT LIKE '%temp%'" before.txt after.txt

# Raw SQL queries
memmondiff --raw-query "SELECT SUM(diff) FROM entries WHERE diff >= 10000" before.txt after.txt
```

## Screenshots

### Linux

![image](https://github.com/user-attachments/assets/fa2e18cc-244a-4979-b5f3-47c49a97773a)

### Windows

- Supported from v0.1.1
- For building from source, use 'CGO_ENABLED=1'. Zig compiler recommended (see justfile).

## Supported Platforms

- Linux
- Windows (requires CGO_ENABLED=1)
- macOS
