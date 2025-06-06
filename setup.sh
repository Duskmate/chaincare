#!/bin/bash

# Variables
REPO_URL="https://github.com/Duskmate/chaincare.git"
TARGET_DIR="$(pwd)"

# Clone the repository
echo "Cloning repository from $REPO_URL..."
git clone $REPO_URL $TARGET_DIR

# Check if clone was successful
if [ $? -ne 0 ]; then
    echo "Failed to clone repository. Exiting..."
    exit 1
fi

echo "Setup complete. Files downloaded."