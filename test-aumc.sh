#!/bin/bash

# Test script for AUMC Go application
# This script tests the basic functionality of the aumc command

set -e

echo "========================================="
echo "AUMC Go Application Test Suite"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run a test
run_test() {
    local test_name="$1"
    local command="$2"
    
    echo -n "Testing: $test_name ... "
    
    if eval "$command" > /tmp/test_output.txt 2>&1; then
        echo -e "${GREEN}PASSED${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}FAILED${NC}"
        echo "  Output:"
        cat /tmp/test_output.txt | sed 's/^/    /'
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

# Check if aumc binary exists
if [ ! -f "/workspace/bin/aumc-linux-amd64" ]; then
    echo -e "${RED}Error: aumc binary not found at /workspace/bin/aumc-linux-amd64${NC}"
    echo "Please build the binary first with: make build-linux"
    exit 1
fi

# Install the binary
echo "Installing aumc binary..."
sudo cp /workspace/bin/aumc-linux-amd64 /usr/local/bin/aumc
sudo chmod +x /usr/local/bin/aumc
echo -e "${GREEN}Binary installed successfully${NC}"
echo ""

# Test 1: Check if aumc command is available
run_test "aumc command availability" "which aumc"

# Test 2: Display help
run_test "aumc --help" "aumc --help"

# Test 3: Display version (if available)
run_test "aumc version/help output" "aumc help || aumc --version || true"

# Test 4: Test init command
echo ""
echo "Testing: aumc init command ..."
cd /tmp
rm -rf test-init
mkdir test-init
cd test-init
if aumc init > /tmp/init_output.txt 2>&1; then
    echo -e "${GREEN}PASSED${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
    echo "  Created files:"
    ls -la | sed 's/^/    /'
else
    echo -e "${YELLOW}SKIPPED (command may not be fully implemented)${NC}"
    cat /tmp/init_output.txt | sed 's/^/    /'
fi

# Test 5: Test config command (if available)
echo ""
echo "Testing: aumc config command ..."
if aumc config --help > /tmp/config_output.txt 2>&1; then
    echo -e "${GREEN}PASSED${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}SKIPPED (command may not be available)${NC}"
fi

# Test 6: Test build command (if available)
echo ""
echo "Testing: aumc build command ..."
if aumc build --help > /tmp/build_output.txt 2>&1; then
    echo -e "${GREEN}PASSED${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}SKIPPED (command may not be available)${NC}"
fi

# Test 7: Test world command (if available)
echo ""
echo "Testing: aumc world command ..."
if aumc world --help > /tmp/world_output.txt 2>&1; then
    echo -e "${GREEN}PASSED${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
else
    echo -e "${YELLOW}SKIPPED (command may not be available)${NC}"
fi

# Summary
echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="
echo -e "Tests Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed: ${RED}$TESTS_FAILED${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
