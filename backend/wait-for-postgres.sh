#!/bin/sh
# wait-for-postgres.sh
# Usage: wait-for-postgres.sh host:port [timeout]

set -e

host_port="$1"
timeout="${2:-60}"

for i in $(seq 1 "$timeout"); do
  if nc -z $(echo "$host_port" | cut -d: -f1) $(echo "$host_port" | cut -d: -f2); then
    echo "Postgres is up!"
    exit 0
  fi
  echo "Waiting for Postgres at $host_port... ($i/$timeout)"
  sleep 1
done

echo "Timeout waiting for Postgres at $host_port"
exit 1
