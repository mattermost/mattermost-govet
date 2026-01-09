#!/bin/bash
# Script to quickly verify we can build and test all analyzers

cd /home/dlauder/Development/mattermost/mattermost-govet

# List analyzers without tests
echo "Analyzers still needing tests:"
for dir in */; do
    if [ -f "${dir}${dir%/}.go" ] && [ ! -f "${dir}${dir%/}_test.go" ]; then
        echo "  - ${dir%/}"
    fi
done

echo ""
echo "Running quick build test..."
go build .
echo "Build: $?"

echo ""
echo "Running all tests..."
go test ./... 2>&1 | grep -E "^(ok|FAIL|\?)" | grep -v "cached"
