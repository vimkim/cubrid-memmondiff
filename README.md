# memmondiff

### Difftool for CUBRID memmon

#### Usage

```bash
Usage: ./memmondiff [options] <before_file> <after_file>

Options:
  --color=MODE               color output (MODE: auto, always, never)
  --sort=TYPE                sort output (TYPE: filename, diff)
```

#### Linux

![image](https://github.com/user-attachments/assets/d9e87217-9eaf-4e69-8d4a-26080f935b4f)

#### Windows

![image](https://github.com/user-attachments/assets/b449799c-515e-43e6-b1ce-2aa5815d00f8)

## Installation

```bash
# Binary
Download the latest release from GitHub Releases

# From source
go install github.com/vimkim/cubrid-memmondiff@latest

# Or
git clone https://github.com/vimkim/cubrid-memmondiff
cd cubrid-memmondiff
go build .

# Or
just run # use casey/just
```
