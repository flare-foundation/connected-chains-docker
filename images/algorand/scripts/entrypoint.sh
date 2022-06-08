#!/bin/bash
set -e

if [ ! -z "${ALGOD_TOKEN}" ]; then
	echo -n "${ALGOD_TOKEN}" > /opt/algorand/.algorand/algod.token
fi

if [[ ! -f /opt/algorand/.algorand/genesis.json ]]; then
	cp /opt/algorand/algorand-defaults/genesis.json /opt/algorand/.algorand/
	if [[ ! -f /opt/algorand/.algorand/config.json ]]; then
		cp /opt/algorand/default-config/config.json /opt/algorand/.algorand/
	fi
	chown -R algo:algo /opt/algorand/.algorand
fi

# Switch to non root
if [ "$(id -u)" = '0' ]; then
	exec gosu algo "$@"
fi

exec "$@"
