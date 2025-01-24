# memmondiff

Difftool for CUBRID memmon that helps analyze memory usage changes with powerful SQL filtering.

## Installation

### Pre-built Binary

Download the latest release from GitHub Releases:

```bash
wget https://github.com/vimkim/cubrid-memmondiff/releases/latest/download/memmondiff-linux-amd64
```

### Using Go

```bash
go install github.com/vimkim/cubrid-memmondiff@latest
```

### Build from Source

```bash
git clone https://github.com/vimkim/cubrid-memmondiff
cd cubrid-memmondiff
go build .
# Or with just: just run
```

## Usage

```text
Usage: ./memmondiff [options] <before_file> <after_file>

Options:
  --color=MODE      color output (MODE: auto(default), always, never)
  (deprecated) --sort=TYPE       sort output (TYPE: diff(default), filename)
  (deprecated) --min=VALUE       minimum diff value to show (default: math.MinInt64)
  (deprecated) --no-new          do not include new entries
  --pretty-print    pretty print numbers
  --sql=FILTER      SQL WHERE clause for filtering (e.g. 'diff >= 10000 AND filename LIKE '%session%')
```

## SQL Filtering Examples

### Show Large Memory Changes

```bash
./memmondiff --sql "diff >= 10000" before.txt after.txt
```

### Filter by File Pattern

```bash
./memmondiff --sql "filename LIKE '%heap%'" before.txt after.txt
```

### Complex Conditions

```bash
./memmondiff --sql "diff >= 5000 AND filename NOT LIKE '%temp%'" before.txt after.txt
```

## Screenshots

### Linux

![image](https://github.com/user-attachments/assets/fa2e18cc-244a-4979-b5f3-47c49a97773a)

### Windows

Coming Soon (used to work for v0.0.6 or before)

## Supported Platforms

- Linux
- Windows
- macOS (theoretically)

```text
The improvements include:
- Added description of SQL filtering capability
- Structured SQL examples section
- Clearer installation options
- Better organization of sections
- Added practical SQL filtering examples
```
