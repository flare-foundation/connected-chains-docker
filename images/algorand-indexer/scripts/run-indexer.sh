#!/bin/bash
set -e

if [ -z "${ALGOD_TOKEN}" ]; then
	ALGOD_TOKEN="$(cat /opt/algorand-indexer/.algorand/algod.token)"
fi

exec algorand-indexer daemon -P "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" --algod-net="${ALGOD_NET}" --algod-token="${ALGOD_TOKEN}"
