# FlowShield Provider

As a service network provider program, miners can deploy network services through this program to obtain rewards.

The program automatically finds through P2P networking and uses PubSub mode to share its own data.

Currently, websocket traffic transmission is supported.

## How does it work?

### Scene Introduction
<img width="1242" alt="image" src="https://user-images.githubusercontent.com/52234994/177236269-03fe1736-66ae-4388-9c3b-3f06f21f3427.png">

## Getting Started

### Configuration

In the configs directory, the user stores the project configuration file with the file name: config.toml, which can be modified according to your own needs

During deployment, you can refer to our current [deployment example](https://github.com/FlowShield/FlowShield/tree/main/deploy/provider) and modify the corresponding configuration


### Quickstart
```shell
$ make
$ ./bin/provider -c ./configs/config.toml
```

## License
FlowShield-Provider uses Apache 2.0 license. See [license](LICENSE) directory for details

## Disclaimers
When you use this software, you have agreed and declared that the author, maintainer and contributor of this software are not responsible for any risks, costs or problems you encounter. If you find a software defect or bug, please submit a patch to help improve!
