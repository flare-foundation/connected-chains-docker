#!/bin/bash
set -e

if [ ! -z "${ALGOD_TOKEN}" ]; then
	echo -n "${ALGOD_TOKEN}" > /opt/algorand/.algorand/algod.token
fi

# Switch to non root
if [ "$(id -u)" = '0' ]; then
	exec gosu algo "$@"
fi

exec "$@"