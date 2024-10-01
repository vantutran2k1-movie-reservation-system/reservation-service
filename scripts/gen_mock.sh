#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <source> <destination>"
    exit 1
fi

SOURCE=$1
DESTINATION=$2

PACKAGE_NAME=$(basename "$(dirname "$DESTINATION")")

mockgen -source="$SOURCE" -destination="$DESTINATION" -package="$PACKAGE_NAME"
