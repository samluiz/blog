#!/bin/bash

set -m

# Function to perform cleanup
cleanup() {
    echo "Cleaning up..."
    # Delete files here
    rm -f tailwindcss tailwindcss.exe
    echo "Cleanup complete."
}

# Trap the SIGINT (Ctrl+C) and SIGTERM signals to execute the cleanup function
trap cleanup EXIT

# Check if argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 [windows|linux|mac]"
    exit 1
fi

# load environment variables
set -a
source .env
set +a

# Check which argument was provided
case "$1" in
    "windows")
        curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-windows-x64.exe
        chmod +x tailwindcss-windows-x64.exe
        mv tailwindcss-windows-x64.exe tailwindcss.exe
        set GOOS=windows
        set GOARCH=amd64
        "C:\Program Files\go\bin\air.exe" -c .air.windows.conf
        ;;
    "linux")
        curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-arm64
        chmod +x tailwindcss-linux-arm64
        mv tailwindcss-linux-arm64 tailwindcss
        export GOOS=linux
        export GOARCH=arm64
        ~/go/bin/air -c .air.linux.conf
        ;;
    "mac")
        curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
        chmod +x tailwindcss-macos-arm64
        mv tailwindcss-macos-arm64 tailwindcss
        export GOOS=darwin
        export GOARCH=arm64
        ~/go/bin/air -c .air.mac.conf
        ;;
    *)
        echo "Invalid argument. Usage: $0 [windows|linux|mac]"
        exit 1
        ;;
esac