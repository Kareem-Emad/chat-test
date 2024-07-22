#!/bin/bash
set -e

# Wait for Cassandra to be ready
until cqlsh -e "describe keyspaces"; do
  >&2 echo "Cassandra is unavailable - sleeping"
  sleep 5
done

>&2 echo "Cassandra is up - executing command"

# Run the CQL script
cqlsh -f /docker-entrypoint-initdb.d/init-cassandra.cql
