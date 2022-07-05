# [Fullnode]

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  [![GoDoc](https://godoc.org/github.com/cloudflare/cfssl?status.svg)](https://github.com/CloudSlit/cloudslit/tree/main/fullnode)

Fullnode is a Dao of [Provider](https://github.com/CloudSlit/cloudslit/tree/main/provider).It shows all Provider node that connects together using P2P networking.And provides users a convenient way to build their own network tunnel.

Every Fullnode needs to deposit a few tokens to become a Dao.After that, other Provider node can join this Dao by PubSub a same topic.Also these nodes will report some metadata including public network IP,wallet,contry and so on.

People can choose any one or two Provider nodes to build their own network tunnel by paying a few tokens.The Provider nodes which been selected will get rewards.

When all things come up,Enjoy CloudSlit.

## Features

- Web3 Dao
- Decentralized storage
- Zero-Trust
- Network proxy

## How does it work



## Building

```shell
$ git clone git@github.com:ztalab/ZAManager.git
$ cd cloudslit/fullnode
$ make
```

You can set GOOS and GOARCH environment variables to allow Go to cross-compile alternative platforms.

The resulting binaries will be in the bin folder:

```shell
$ tree bin
bin
├── fullnode
```

## Installing

### Docker-compose

~~~shell
cd deploy/docker-compose
vim nginx.conf. line 38,chane you own domain
docker-compose up -d
~~~

This will also install [CA](https://github.com/CloudSlit/cloudslit/tree/main/ca) and [Portal](https://github.com/CloudSlit/cloudslit/tree/main/portal)

Don't forget change  `zta_oauth2 ` table Oauth2 data in mysql with your owns

## License

Fullnode source code is available under the Apache 2.0 [License](https://github.com/CloudSlit/cloudslit/blob/main/fullnode/LICENSE).
