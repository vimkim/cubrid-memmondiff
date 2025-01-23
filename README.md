# memmondiff

### Difftool for CUBRID memmon

#### Linux

![image](https://github.com/user-attachments/assets/7ce442aa-f179-4443-b9d8-865ac37d67da)

#### Windows

![image](https://github.com/user-attachments/assets/dfce602e-a012-435f-9528-3abf2bc430cc)

## Installation

```bash
# Binary
Download the latest release from GitHub Releases

# From source
go install github.com/vimkim/cubrid-memmondiff@latest
```

## Usage

```bash
memmondiff [--color=always|auto|never] file1.txt file2.txt
```

Options:

- `--color=always`: Always show colors
- `--color=auto`: Show colors when output is to terminal (default)
- `--color=never`: Never show colors

Example output:

```
base/system_parameter.c:11736 | 1032 (=2064-1032)
base/system_parameter.c:11927 | 128 (=256-128)
```

Colors:

- Red: Memory increase
- Green: Memory decrease
