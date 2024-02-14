#!/bin/bash

# Define install location
INSTALL_DIR=/usr/local/bin
BINARY_NAME=timer

# macOS-specific: Check for Homebrew and suggest installation if not found
if ! command -v brew &> /dev/null
then
    echo "Homebrew not found. Consider installing Homebrew for managing packages on macOS."
    echo "Visit https://brew.sh/ for instructions."
    exit 1
fi

# Check if already installed and remove
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    echo "Existing installation found. Removing..."
    sudo rm "$INSTALL_DIR/$BINARY_NAME"
fi

# Copy binary to install location
echo "Installing $BINARY_NAME to $INSTALL_DIR"
sudo cp "./$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

# Make binary executable
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "$BINARY_NAME installed successfully."
