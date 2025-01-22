#!/bin/bash
set -euo pipefail

TEST_NAME=$(basename "$0")
TEST_DIR="tests/testdata"
BINARY="memmondiff"
ACTUAL="${TEST_DIR}/actual.txt"
EXPECTED="${TEST_DIR}/expected_1_2.txt"

echo "Running test: ${TEST_NAME}"

# Build and verify
if ! go build -o ${BINARY} .; then
    echo "FAIL: Build failed"
    exit 1
fi

# Run test
if ! ./${BINARY} --color=never ${TEST_DIR}/mem1.txt ${TEST_DIR}/mem2.txt >${ACTUAL}; then
    echo "FAIL: Program execution failed"
    rm -f ${ACTUAL} ${BINARY}
    exit 1
fi

# Compare results
if diff "${EXPECTED}" "${ACTUAL}" >/dev/null; then
    echo "PASS: ${TEST_NAME}"
else
    echo "FAIL: ${TEST_NAME}"
    echo "Differences found:"
    diff -u "${EXPECTED}" "${ACTUAL}"
    delta "${EXPECTED}" "${ACTUAL}"
fi

# Cleanup
rm -f ${ACTUAL} ${BINARY}
