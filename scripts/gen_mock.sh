#!/bin/bash

# Check if the number of arguments is less than 1
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <source>"
    exit 1
fi

# Assign the source argument to a variable
SOURCE=$1

# Extract folder name and file name from source path
FOLDER_NAME=$(basename "$(dirname "$SOURCE")")
FILE_NAME=$(basename "$SOURCE")

# Construct the destination path
MOCK_DESTINATION="app/mocks/mock_${FOLDER_NAME}/${FILE_NAME}"

# Create the destination directory if it doesn't exist
mkdir -p "$(dirname "$MOCK_DESTINATION")"

# Extract package name from the destination folder
PACKAGE_NAME=$(basename "$(dirname "$MOCK_DESTINATION")")

# Run mockgen command with dynamic destination and package name
mockgen -source="$SOURCE" -destination="$MOCK_DESTINATION" -package="$PACKAGE_NAME"

echo "Mock generated at $MOCK_DESTINATION"
