#!/bin/bash

# Bor
docker pull registry.hub.docker.com/0xpolygon/bor:0.3.0 
docker tag registry.hub.docker.com/0xpolygon/bor:0.3.0 registry.hub.docker.com/flarefoundation/bor:0.3.0
docker push registry.hub.docker.com/flarefoundation/bor:0.3.0

# Heimdall
docker pull registry.hub.docker.com/0xpolygon/heimdall:0.3.0
docker tag registry.hub.docker.com/0xpolygon/heimdall:0.3.0 registry.hub.docker.com/flarefoundation/heimdall:0.3.0
docker push registry.hub.docker.com/flarefoundation/heimdall:0.3.0
