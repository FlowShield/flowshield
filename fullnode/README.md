# Fullnode

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  [![GoDoc](https://godoc.org/github.com/cloudflare/cfssl?status.svg)](https://github.com/FlowShield/FlowShield/tree/main/fullnode)

Fullnode is a Dao of [Provider](https://github.com/FlowShield/FlowShield/tree/main/provider).It shows all Provider node that connects together using P2P networking.And provides users a convenient way to build their own network tunnel.

Every Fullnode needs to deposit a few tokens to become a Dao.After that, other Provider node can join this Dao by PubSub a same topic.Also these nodes will report some metadata including public network IP,wallet,contry and so on.

People can choose any one or two Provider nodes to build their own network tunnel by paying a few tokens.The Provider nodes which been selected will get rewards.

When all things come up,Enjoy FlowShield.

## Features

- Web3 Dao
- Decentralized storage
- Zero-Trust
- Network Security

## How does it work



## Installing

### Docker-compose

1. Get deploy Dockerfile

~~~shell
cd deploy/docker-compose
~~~

just pay attention to `nginx.conf` and `docker-compose.yaml`

2. Change `nginx.conf` line32 to your own domain

~~~nginx
...
listen 80 default_server;

server_name dash.FlowShield.xyz;  #here is your domain

root /usr/share/nginx/html;
index index.html;
...
~~~

3. Change some ENV in `docker-compose.yaml` file

~~~ini
CS_OAUTH2_CLIENT_ID: 'your client id should apply from github'
CS_OAUTH2_CLIENT_SECRET: 'your client secret should apply from github'
CS_PRIVATE_KEY: 'a blockchain private key'
CS_CONTRACT_TOKEN: 'contract address'
CS_W3S_TOKEN: 'apply from https://web3.storage'
~~~

This will also install [CA](https://github.com/FlowShield/FlowShield/tree/main/ca) and [Portal](https://github.com/FlowShield/FlowShield/tree/main/portal)

## License

Fullnode source code is available under the Apache 2.0 [License](https://github.com/FlowShield/FlowShield/blob/main/fullnode/LICENSE).
