# Set up variables
BINARY := "lbox"
SRC := "./cmd/main.go"
BUILD_DIR := "./bin"

# Command to run the app
run:
    go run {{SRC}}

# Build the binary
build:
    mkdir -p {{BUILD_DIR}}
    go build -o {{BUILD_DIR}}/{{BINARY}} {{SRC}}

# Clean the build directory
clean:
    rm -rf {{BUILD_DIR}}
