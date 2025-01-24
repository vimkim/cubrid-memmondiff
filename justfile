set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]

build:
    go build .

run: run-ycsb

build-release-from-linux:
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o memmondiff-linux-amd64
    CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o memmondiff-windows-amd64.exe

run-ycsb: build
    ./memmondiff ./testdata/ycsb_before ./testdata/ycsb_after

run-min-10000: build
    ./memmondiff --min=10000 ./testdata/ycsb_before ./testdata/ycsb_after

run-min-10: build
    ./memmondiff --min=10 ./testdata/ycsb_before ./testdata/ycsb_after

run-no-new: build
    ./memmondiff --no-new --min=10000 ./testdata/ycsb_before ./testdata/ycsb_after

tests: test-ycsb

test-ycsb: build
    ./test/test_ycsb.sh

help: build
    ./memmondiff --help

pretty-print: build
    ./memmondiff --min=10000 --pretty-print ./testdata/ycsb_before ./testdata/ycsb_after

run-sql-min-10000-no-new: build
    ./memmondiff --sql="diff >= 10000 AND before <> 0" ./testdata/ycsb_before ./testdata/ycsb_after
