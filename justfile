build:
    go build .

run: run-ycsb

run-example:
    go run . 1.txt 2.txt

build-release:
    GOOS=linux GOARCH=amd64 go build -o memmondiff-linux-amd64
    GOOS=windows GOARCH=amd64 go build -o memmondiff-windows-amd64.exe

run-ycsb:
    go run . tests/testdata/ycsb_before tests/testdata/ycsb_after

tests: test-ycsb

test-example: build
    ./tests/test_1_2.sh

test-ycsb: build
    ./tests/test_ycsb.sh

help: build
    ./memmondiff --help
