#!/bin/bash

# Define install location and binary name
INSTALL_DIR=/usr/local/bin
BINARY_NAME=timer

# Remove binary
if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
    echo "Removing $BINARY_NAME..."
    sudo rm "$INSTALL_DIR/$BINARY_NAME"
    echo "$BINARY_NAME removed successfully."
else
    echo "$BINARY_NAME is not installed."
fi
