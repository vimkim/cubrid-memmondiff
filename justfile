build:
    go build .

run: run-ycsb

run-example:
    ./memmondiff 1.txt 2.txt

build-release:
    GOOS=linux GOARCH=amd64 go build -o memmondiff-linux-amd64
    GOOS=windows GOARCH=amd64 go build -o memmondiff-windows-amd64.exe

run-ycsb: build
    ./memmondiff ./testdata/ycsb_before ./testdata/ycsb_after

tests: test-ycsb

test-example: build
    ./test/test_1_2.sh

test-ycsb: build
    ./test/test_ycsb.sh

help: build
    ./memmondiff --help
