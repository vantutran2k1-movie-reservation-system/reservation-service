#!/bin/bash
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-5432}
export DB_USER=${DB_USER:-postgres}
export DB_PASSWORD=${DB_PASSWORD:-postgres}
export DB_NAME=${DB_NAME:-booking}

MIGRATE_COMMAND="migrate -path migrations -database \"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable\" up"

get_second_latest_version() {
  second_latest_version=$(ls migrations | awk -F'_' '{print $1}' | sort -n | uniq | tail -2 | head -1)
  echo "$second_latest_version"
}

if [[ "$1" == "--force" ]]; then
  echo "Running migrations with --force..."

    VERSION=$(get_second_latest_version)

    if [ -z "$VERSION" ]; then
      echo "No valid migrations found in the folder."
      exit 1
    fi

    SQL_QUERY="UPDATE schema_migrations SET version = '$VERSION', dirty = false;"

    echo "Updating schema_migrations table..."
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_QUERY"

    if [ $? -eq 0 ]; then
      echo "Schema migrations table updated successfully."

      echo "Running migrations..."
      eval "$MIGRATE_COMMAND"
    else
      echo "Failed to update the schema migrations table."
      exit 1
    fi
else
  eval "$MIGRATE_COMMAND"
fi

