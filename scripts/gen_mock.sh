#!/bin/bash

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <source> <destination> <package>"
    exit 1
fi

SOURCE=$1
DESTINATION=$2

PACKAGE_NAME=$(basename "$(dirname "$DESTINATION")")

mockgen -source="$SOURCE" -destination="$DESTINATION" -package="$PACKAGE_NAME"
