#!/bin/bash

# Create local Go workspace
mkdir -p "$HOME/go-workspace/pkg/mod"
mkdir -p "$HOME/go-workspace/bin"

# Set environment variables for this session
export GOPATH="$HOME/go-workspace"
export GOMODCACHE="$HOME/go-workspace/pkg/mod"
export GOBIN="$HOME/go-workspace/bin"

# Show confirmation
echo "Go environment configured:"
echo "GOPATH=$GOPATH"
echo "GOMODCACHE=$GOMODCACHE"
echo "GOBIN=$GOBIN"

# Optional: initialize module if not already inside one
if [ ! -f "go.mod" ]; then
  echo "No go.mod found. Initializing module..."
  go mod init myproject
fi

# Download dependency
go get github.com/mattn/go-sqlite3

echo "Done."