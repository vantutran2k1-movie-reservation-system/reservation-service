#!/bin/bash

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <migration_name>"
  exit 1
fi

MIGRATION_NAME=$1
TIMESTAMP=$(date -u +"%Y%m%d%H%M%S")
MIGRATION_DIR="./migrations"
UP_FILE="${MIGRATION_DIR}/${TIMESTAMP}_${MIGRATION_NAME}.up.sql"
DOWN_FILE="${MIGRATION_DIR}/${TIMESTAMP}_${MIGRATION_NAME}.down.sql"

mkdir -p $MIGRATION_DIR

touch $UP_FILE
touch $DOWN_FILE

echo "Migration files created:"
echo " - $UP_FILE"
echo " - $DOWN_FILE"
