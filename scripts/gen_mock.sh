#!/bin/bash

# Check if the number of arguments is less than 2
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <source> <destination> <package>"
    exit 1
fi

# Assign arguments to variables
SOURCE=$1
DESTINATION=$2
PACKAGE=$3

# Run mockgen command
mockgen -source="$SOURCE" -destination="$DESTINATION" -package="$PACKAGE"
