# memmondiff

## Difftool for CUBRID memmon

### Usage

```txt
Usage: ./memmondiff [options] <before_file> <after_file>

Options:
  --color=MODE      color output (MODE: auto, always, never)
  --sort=TYPE       sort output (TYPE: filename, diff)
  --min=VALUE       minimum diff value to show (default: show all)
  --no-new          do not include new entries
  --pretty-print    pretty print numbers
```

### Installation

```bash
# Binary
# Just download the latest release from GitHub Releases
# For example,
wget https://github.com/vimkim/cubrid-memmondiff/releases/download/v0.0.2/memmondiff-linux-amd64

###

# Or use go module,
go install github.com/vimkim/cubrid-memmondiff@latest

###

# Or build from source,
git clone https://github.com/vimkim/cubrid-memmondiff
cd cubrid-memmondiff
go build .

###

# You can use Casey/just
just run # use casey/just
```

#### Filtering Examples

##### No New Entries, and Only Diff >= 10000

![image](https://github.com/user-attachments/assets/fa2e18cc-244a-4979-b5f3-47c49a97773a)

### Supported Platforms

- Linux
- Windows
- MacOs (theoretically)

#### Linux

![image](https://github.com/user-attachments/assets/d9e87217-9eaf-4e69-8d4a-26080f935b4f)

#### Windows

![image](https://github.com/user-attachments/assets/b449799c-515e-43e6-b1ce-2aa5815d00f8)
