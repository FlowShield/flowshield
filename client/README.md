# FlowShield Client

This is the FlowShield client program.

You need to communicate with the whole node through this program to obtain the configuration and zero trust edge calculation rules. Communicate with network providers through this program.

## How does it work?

### Scene Introduction
<img width="1242" alt="image" src="https://user-images.githubusercontent.com/52234994/177236269-03fe1736-66ae-4388-9c3b-3f06f21f3427.png">

## Getting Started

### Configuration

In the configs directory, the user stores the project configuration file with the file name: config.toml, which can be modified according to your own needs

By default, you do not need to modify the configuration file, just refer to QuickStart.

### Quickstart
```shell
$ make
$ ./bin/client
```

## License
FlowShield-client uses Apache 2.0 license. See [license](LICENSE) directory for details

## Disclaimers
When you use this software, you have agreed and declared that the author, maintainer and contributor of this software are not responsible for any risks, costs or problems you encounter. If you find a software defect or bug, please submit a patch to help improve!
