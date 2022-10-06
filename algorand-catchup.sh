#!/bin/bash
set -eu

ALGOD_CONTAINER=$(docker-compose ps -q algorand)

ALGOD_PORT=$(docker port ${ALGOD_CONTAINER} 8080/tcp | grep "\." | cut -d ":" -f 2)
ALGOD_URL="http://127.0.0.1:${ALGOD_PORT}"

until curl -f -m 3 "${ALGOD_URL}/health" >/dev/null 2>&1; do
	echo "waiting for algod"
	sleep 3
done

sleep 1

ALGOD_ADMIN_TOKEN=$(docker exec ${ALGOD_CONTAINER} cat /opt/algorand/.algorand/algod.admin.token)

LAST_ROUND=$(curl -s -H "X-Algo-API-Token: ${ALGOD_ADMIN_TOKEN}" "${ALGOD_URL}/v2/status" | jq -r '.["last-round"]')
CURRENT_CATCHPOINT=$(curl -s -H "X-Algo-API-Token: ${ALGOD_ADMIN_TOKEN}" "${ALGOD_URL}/v2/status" | jq -r '.catchpoint')

ALGOD_NETWORK=$(docker exec ${ALGOD_CONTAINER} printenv ALGOD_NETWORK) || true

if [ -z "${ALGOD_NETWORK}" ]; then
	ALGOD_NETWORK="mainnet"
fi

echo "last round: ${LAST_ROUND}"
echo "current checkpoint: ${CURRENT_CATCHPOINT}"

if [ "$LAST_ROUND" -lt "20000000" ] && [[ -z "$CURRENT_CATCHPOINT" ]]; then
	LATEST_CATCHPOINT=$(curl -s "https://algorand-catchpoints.s3.us-east-2.amazonaws.com/channel/${ALGOD_NETWORK}/latest.catchpoint")
	echo "catching up to catchpoint ${LATEST_CATCHPOINT}"
	LATEST_CATCHPOINT=$(echo "${LATEST_CATCHPOINT}" | sed "s/#/%23/g")
	curl -s -X POST -H "X-Algo-API-Token: ${ALGOD_ADMIN_TOKEN}" "${ALGOD_URL}/v2/catchup/${LATEST_CATCHPOINT}"
else
	echo "already catching up or previously initialized"
fi
