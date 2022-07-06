#!/bin/bash

exec > >(awk '{ print strftime("[%Y-%m-%d %H:%M:%S]"), $0 }') 2>&1

RPC_PASS="${RPC_PASS:-123456}"
HOST="${HOST:-dogecoind}"
PORT="${PORT:-22555}"

trap trapint SIGINT SIGTERM
function trapint {
    echo "Received SIGINT or SIGTERM, exiting"
    exit 0
}

function hc() {
    OUT=$(curl -vvv -s -X POST -m 10 -H "Content-type: application/json" -d '{"jsonrpc": "1.0", "id":"hc", "method": "getblockchaininfo", "params":[]}' -u admin:$RPC_PASS http://$HOST:$PORT)

    if [ $? -ne 0 ]; then
        echo "503"
    else
        jq -e '.result.initialblockdownload == false' <<< "$OUT" > /dev/null
    
        if [ $? -eq 0 ]; then
            echo "200"
        else
            echo "503"
        fi
    fi
}

while [ true ]
do
	echo -e "HTTP/1.1 $(hc)\r\nContent-Length: 0\r\n" |  nc -vl 11111 | ts '[%Y-%m-%d %H:%M:%S]'
	sleep 0.05
done