# Create a local cache in your home directory
mkdir -p ~/.cache/go-mod

# Set it for this session
export GOMODCACHE="$HOME/.cache/go-mod"

# Make it permanent for future sessions
echo 'export GOMODCACHE="$HOME/.cache/go-mod"' >> ~/.bashrc

# Run again
go run main.go