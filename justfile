build:
    go build .

run:
    go run . 1.txt 2.txt

build-release:
    GOOS=linux GOARCH=amd64 go build -o memmondiff-linux-amd64
    GOOS=windows GOARCH=amd64 go build -o memmondiff-windows-amd64.exe

tests:
    ./tests/test_1_2.sh
