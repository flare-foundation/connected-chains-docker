#!/bin/bash

PASS="${1:-}"

if [[ ! -f jwt.hex ]]; then
    secret=$PASS
    if [[ ${#PASS} -eq 0 ]]; then
	    openssl rand -hex 32 | tr -d "\n" > jwt.hex
    else
        echo $secret > jwt.hex
    fi
    echo "Your generated jwt is: $secret"
else
    echo "jwt already exists, skipping secret generation"
fi