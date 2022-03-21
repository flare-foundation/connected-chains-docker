#!/bin/bash
set -e

# Switch to non root
if [ "$(id -u)" = '0' ]; then
	exec gosu algoindexer "$@"
fi

exec "$@"