#!/bin/bash

echo "Building binaries"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/emperium-x86_64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/emperium-aarch64
echo "Looking for elf compressor"
which upx > /dev/null
if [ $? -eq 1 ]; then
    echo "Not found, binaries wont be compressed"
else 
    upx ./bin/emperium-x86_64
    upx ./bin/emperium-aarch64
fi