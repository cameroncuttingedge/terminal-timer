#!/bin/bash

# Define install location
INSTALL_DIR=/usr/local/bin
BINARY_NAME=timer

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
