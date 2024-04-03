#!/bin/bash

docker -v
if [ $? -ne 0 ]; then
    echo "Installing docker"
    sudo apt-get -y update
    sudo apt-get -y install ca-certificates curl gnupg lsb-release jq
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get -y update
    sudo apt-get -y install docker-ce docker-ce-cli containerd.io
fi    

# openssl
openssl version
if [[ $? -ne 0 ]]; then
    echo "Installing openssl"
    sudo apt-get -y install openssl
fi

# jq
jq --version
if [[ $? -ne 0 ]]; then
    echo "Installing jq"
    sudo apt-get -y install jq
fi

# rpcauth
NETWORK="${1:-mainnet}"
PASS="${2:-}"

CONFIG_DIR="config"

if [[ "$NETWORK" = "testnet" ]]; then
    CONFIG_DIR="config-testnet"
fi

cd "./${CONFIG_DIR}/bitcoin"
./rpcauth.py admin $PASS

cd ../litecoin
./rpcauth.py admin $PASS

cd ../dogecoin
./rpcuser.py admin $PASS

cd ../ripple
#./rpcauth.py admin $PASS

cd ../algorand
bash gen_auth_token.sh $PASS

cd ../..

echo "Installation finished."
