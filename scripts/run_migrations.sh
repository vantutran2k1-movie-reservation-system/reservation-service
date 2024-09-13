#!/bin/bash
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-5432}
export DB_USER=${DB_USER:-postgres}
export DB_PASSWORD=${DB_PASSWORD:-postgres}
export DB_NAME=${DB_NAME:-booking}

MIGRATE_COMMAND="migrate -path migrations -database \"postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable\" up"

# Function to get the second latest migration version
get_second_latest_version() {
  # List the migration files, extract timestamps, sort them, and select the second latest
  second_latest_version=$(ls migrations | awk -F'_' '{print $1}' | sort -n | uniq | tail -2 | head -1)
  echo "$second_latest_version"
}

# Check if the script is run with --force
if [[ "$1" == "--force" ]]; then
  echo "Running migrations with --force..."

  # Get the second latest migration version
    VERSION=$(get_second_latest_version)

    # Check if a version was found
    if [ -z "$VERSION" ]; then
      echo "No valid migrations found in the folder."
      exit 1
    fi

    # Update the schema_migrations table
    SQL_QUERY="UPDATE schema_migrations SET version = '$VERSION', dirty = false;"

    # Execute the SQL query using psql
    echo "Updating schema_migrations table..."
    # shellcheck disable=SC2086
    PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$SQL_QUERY"

    # Check if the update was successful
    if [ $? -eq 0 ]; then
      echo "Schema migrations table updated successfully."

      # Run the migration after the update
      echo "Running migrations..."
      eval "$MIGRATE_COMMAND"
    else
      echo "Failed to update the schema migrations table."
      exit 1
    fi
else
  eval "$MIGRATE_COMMAND"
fi

