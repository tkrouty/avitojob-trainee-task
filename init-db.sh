#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE accounts(
    user_id VARCHAR UNIQUE NOT NULL,
    balance DECIMAL NOT NULL
    );
    CREATE TABLE transaction_history(
    transaction_id SERIAL PRIMARY KEY,
    source_id VARCHAR NOT NULL,
    target_id VARCHAR NOT NULL,
    sum DECIMAL NOT NULL,
    transaction_time TIMESTAMP NOT NULL
    );
EOSQL
