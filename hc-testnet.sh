#!/bin/bash

PASS="${1:-admin}"

echo "====================== BITCOIND ======================"
curl -X POST -m 10 -H "Content-type: application/json" -d '{"jsonrpc": "1.0", "id":"hc", "method": "getblockchaininfo", "params":[]}' http://admin:$PASS@localhost:18332 |jq

echo "====================== LITECOIND ======================"
curl -X POST -m 10 -H "Content-type: application/json" -d '{"jsonrpc": "1.0", "id":"hc", "method": "getblockchaininfo", "params":[]}' http://admin:$PASS@localhost:19332 |jq

echo "====================== DOGECOIND ======================"
curl -X POST -m 10 -H "Content-type: application/json" -d '{"jsonrpc": "1.0", "id":"hc", "method": "getblockchaininfo", "params":[]}' http://admin:$PASS@localhost:44555 |jq

echo "====================== RIPPLED ======================"
curl -X POST -m 10 -H "Content-type: application/json" -d '{"method": "server_info", "params":[{"api_version": 1}]}' http://localhost:11234 |jq
