#!/bin/bash
# Runs automatically on first Postgres container start (docker-entrypoint-initdb.d
# convention). Creates one database per service — separate DBs, not separate
# schemas, so each service's data is fully isolated (closer to how you'd run
# this in real microservices, even though it's one Postgres instance for now).
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE DATABASE authdb;
    CREATE DATABASE orderdb;
EOSQL
