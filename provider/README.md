# [Provider]

As a service network provider program, miners can deploy network services through this program to obtain rewards.

The program automatically finds through P2P networking and uses PubSub mode to share its own data.

Currently, websocket traffic transmission is supported.

## How does it work?

### Scene Introduction
<img width="1242" alt="image" src="https://user-images.githubusercontent.com/52234994/177236269-03fe1736-66ae-4388-9c3b-3f06f21f3427.png">

## Getting Started

### Configuration

In the config directory, the user stores the project configuration file with the file name: config server Yaml, which can be modified according to your own needs

You only need to modify the following contents, and other configurations can be modified according to your own needs.

```yaml
# Information about themselves
[Common]
  # Self information
  UniqueId = "provider"
  AppName = "provider"
  # Wallet address
  PeerId = "0x1B4b827703dc3545089fcee70F0e6e732BFF4413"
  # External service address, Need to support websocket protocol
  LocalAddr = "server.com"
  # External service port
  LocalPort = 5091
```

### Quickstart
```shell
$ make
$ ./bin/backend -c ./configs/config_server.yaml
```

## License
Icefiredb agent uses Apache 2.0 license. See [license](.License) directory for details

## Disclaimers
When you use this software, you have agreed and declared that the author, maintainer and contributor of this software are not responsible for any risks, costs or problems you encounter. If you find a software defect or bug, please submit a patch to help improve!
