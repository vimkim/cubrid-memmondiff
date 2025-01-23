#!/bin/bash
set -euo pipefail

TEST_NAME=$(basename "$0")
TEST_DIR="tests/testdata"
BINARY="memmondiff"
ACTUAL="${TEST_DIR}/actual_ycsb.txt"
EXPECTED="${TEST_DIR}/expected_ycsb.txt"

echo "Running test: ${TEST_NAME}"

# Run test
if ! ./${BINARY} --color=never ${TEST_DIR}/ycsb_before ${TEST_DIR}/ycsb_after >${ACTUAL}; then
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
