build:
    go build .

run: run-ycsb

build-release:
    GOOS=linux GOARCH=amd64 go build -o memmondiff-linux-amd64
    GOOS=windows GOARCH=amd64 go build -o memmondiff-windows-amd64.exe

run-ycsb: build
    ./memmondiff ./testdata/ycsb_before ./testdata/ycsb_after

run-min: build
    ./memmondiff --min=10000 ./testdata/ycsb_before ./testdata/ycsb_after

tests: test-ycsb

test-ycsb: build
    ./test/test_ycsb.sh

help: build
    ./memmondiff --help
