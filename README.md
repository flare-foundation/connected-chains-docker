# Flare/Songbird Attestation Provider - Connected Chains Docker images

A quickstart repo filled with Docker images which allows people to get up and running with the chains connected to the Flare Network.

The following nodes are included:
- [Bitcoin](https://github.com/bitcoin/bitcoin)
- [Litecoin](https://github.com/litecoin-project/litecoin)
- [Dogecoin](https://github.com/dogecoin/dogecoin)
- [Rippled](https://github.com/ripple/rippled)
- [Algorand](https://github.com/algorand/go-algorand)

You can use this repo to get you started but we encourage the community to configure their own setups for increased technical diversity of the attestation provider ecosystem.

# System Requirements

The following specifications were observed to be able to run all nodes on a single machine. You probably want to run one node per system in production for better performance and uptime.

- OS: Ubuntu 20.04
- CPU Cores: 16
- RAM: 64GB
- Disk space: 3TB with an option or plan in place to expand capacity. SSD is recommended, some chains like xrpl can fall out of sync on regular disks.

Bootstrap time depends on your infrastructure and network, in our testing it is a few hours for litecoin, dogecoin and rippled, more than a day for bitcoin and weeks for algorand.

# Installation
```
git clone https://github.com/flare-foundation/connected-chains-docker /opt/connected-chains
cd /opt/connected-chains
./install.sh <your-provided-password>
```

# Running

All containers:
```
cd /opt/connected-chains
docker-compose up -d
```

Single container:
```
docker-compose up bitcoin
```

Stop a single container:
```
docker-compose stop bitcoin
```

You can check the bootstrap process with the `hc.sh` script. `hc <your-provided-password>`

# Logs

```
docker-compose logs -f --tail=1000 bitcoin
```

# Node configuration
Each node has a configuration file provided in `/opt/connected-chains/<node>/<config-file-name>.conf`.
Config files are volume-mounted to each container. After changing the config file you should restart the container.

# Node data
Data for each node is stored in a docker volume. To find the exact mount point on your filesystem, run
```
docker volume ls
sudo docker inspect <volume-name> | grep Mountpoint
```

# Security considerations
Installation script will create a username and password for nodes and insert a line into each config file.
You should save the usernames and passwords in your password manager for later use. You can always change them by editing the config files manually.

Ports bind to `127.0.0.1` by default so nothing is exposed to the public internet. Change the `BIND_IP` variable in `.env` file to bind to other interfaces.

If you are running attestation clients on a different machine, consider locking down RPC port access to the specific IP with firewall.

If you are running the clients in different networks you might also want to consider running TLS, either natively on each node that supports it or behind a reverse proxy.

# Mounted storage 
If you mount additional storage and want it to be used by the chains, change docker data to your mounted directory:

In `/etc/docker/daemon.json` add/change:
```
{
  "data-root": "/data"
}
```

followed by `sudo systemctl restart docker`.

Alternatively, if you do not wish to change data directory for your docker daemon you can switch to bind volume mounts or volume mounts with nfs driver in the compose file.

# Building your own images

All Dockerfile definitions are in `images` folder. Images use Moby BuildKit extensions.

Common build problems:
- OOM kills compiler. Increase memory of your docker daemon (if using Docker desktop), increase Docker daemon swap and lower the parallel jobs flag (the `-j X` parameter).
- if compiler is killed by OOM, image can continue building as if nothing bad happened. Makes sure to check the logs to confirm that build process succeeded. Otherwise container won't start due to missing binaries.
- git clone hangs: try to increase the setting `git config --global http.postBuffer <max-bytes>`

# License

License: MIT
Copyright 2022 Flare Foundation
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.