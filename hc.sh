#!/bin/bash

PASS="${1:-admin}"

echo "====================== BITCOIND ======================"
curl -X POST -m 10 -H "Content-type: application/json" -d '{"jsonrpc": "1.0", "id":"hc", "method": "getblockchaininfo", "params":[]}' http://admin:$PASS@localhost:9332 |jq

echo "====================== LITECOIND ======================"
curl -X POST -m 10 -H "Content-type: application/json" -d '{"jsonrpc": "1.0", "id":"hc", "method": "getblockchaininfo", "params":[]}' http://admin:$PASS@localhost:8332 |jq

echo "====================== DOGECOIND ======================"
curl -X POST -m 10 -H "Content-type: application/json" -d '{"jsonrpc": "1.0", "id":"hc", "method": "getblockchaininfo", "params":[]}' http://admin:$PASS@localhost:7332 |jq

echo "====================== RIPPLED ======================"
curl -X POST -m 10 -H "Content-type: application/json" -d '{"method": "server_info", "params":[{"api_version": 1}]}' http://localhost:51234 |jq

echo "====================== ALGORAND ======================"
curl -X GET -m 10 http://localhost:6332/v2/status -H "X-Algo-API-Token: $PASS" |jq
