#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <source>"
    exit 1
fi

SOURCE=$1
EXCLUDE_SUFFIXES="_test,_mock"
IFS=',' read -ra SUFFIXES <<< "$EXCLUDE_SUFFIXES"

generate_mock() {
    local src_file
    local folder_name
    local file_name
    local mock_destination
    local package_name

    src_file=$1
    folder_name=$(basename "$(dirname "$src_file")")
    file_name=$(basename "$src_file")
    mock_destination="app/mocks/mock_${folder_name}/${file_name}"
    package_name=$(basename "$(dirname "$mock_destination")")

    mkdir -p "$(dirname "$mock_destination")"

    mockgen -source="$src_file" -destination="$mock_destination" -package="$package_name"
    echo "Mock generated at $mock_destination"
}

should_exclude_file() {
    local file suffix
    file="$1"

    for suffix in "${SUFFIXES[@]}"; do
        if [[ "$file" == *"$suffix.go" ]]; then
            return 0
        fi
    done
    return 1
}

if [ -f "$SOURCE" ]; then
    if ! should_exclude_file "$SOURCE"; then
        generate_mock "$SOURCE"
    else
        echo "Skipping excluded file: $SOURCE"
    fi
elif [ -d "$SOURCE" ]; then
    for file in "$SOURCE"/*.go; do
        if [ -f "$file" ] && ! should_exclude_file "$file"; then
            generate_mock "$file"
        else
            echo "Skipping excluded file: $file"
        fi
    done
else
    echo "Invalid source path: $SOURCE"
    exit 1
fi