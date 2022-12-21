#!/bin/bash

# Geth
docker pull registry.hub.docker.com/ethereum/client-go:stable
docker tag registry.hub.docker.com/ethereum/client-go:stable registry.hub.docker.com/flarefoundation/geth:stable
docker push registry.hub.docker.com/flarefoundation/geth:stable

# Prysm
docker pull gcr.io/prysmaticlabs/prysm/beacon-chain:stable
docker tag gcr.io/prysmaticlabs/prysm/beacon-chain:stable registry.hub.docker.com/flarefoundation/prysm:stable
docker push registry.hub.docker.com/flarefoundation/prysm:stable
