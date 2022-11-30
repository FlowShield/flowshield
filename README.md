<p align="center">
<img width="800" alt="image" src="https://user-images.githubusercontent.com/34047788/204796156-c5c9d228-725b-4e14-b8ed-1e705ba19bc5.png">

</p>




# FlowShield - Private retrieval of data
- [FlowShield - Private retrieval of data](#flowshield---private-retrieval-of-data)
- [inspiration](#inspiration)
- [Its value](#its-value)
- [How do we build it?](#how-do-we-build-it)
  - [Part one:FlowShield-Fullnode(Ful nodes of private data retrieval network based on DAO Tools)](#part-oneflowshield-fullnodeful-nodes-of-private-data-retrieval-network-based-on-dao-tools)
  - [Part two: FlowShield-Provider(Network Miner, a Secure Network Tunnel Provider for Decentralized Data Private Retrieval)](#part-two-flowshield-providernetwork-miner-a-secure-network-tunnel-provider-for-decentralized-data-private-retrieval)
  - [Part three:FlowShield-Contracts(support:EVM Chains)](#part-threeflowshield-contractssupportevm-chains)
  - [Part four:FlowShield-verifier(Decentralized network quality checker)](#part-fourflowshield-verifierdecentralized-network-quality-checker)
  - [Part five:DeCA(Decentralized PKI CA center)](#part-fivedecadecentralized-pki-ca-center)
  - [Part six:FlowShield-Client(client for private data retrieval.)](#part-sixflowshield-clientclient-for-private-data-retrieval)
- [Disclaimers](#disclaimers)
- [Thanks supports](#thanks-supports)


# inspiration

At present, the options available for interactive (low-latency) communication with privacy guarantee are very limited, and the solutions developed so far all focus on the traditional web model of single source data publisher, and it has defects in delay and threat models.

FlowShield uses blockchain, web3 and secure network technology of private data retrieval to enhance and improve network security/privacy protection of users' privatization.

In order to protect the public's network security under web2, a very popular zero-trust security architecture has emerged. Our team has been working on open source products with zero trust security, but we found that although many zero trust network security companies provide zero trust security platforms, they monopolize users' network access nodes and centrally store users' core security profiles. Therefore, we are considering whether we can use web3 technology to realize a secure network for private data retrieval. We designed FlowShield project to provide users with a decentralized secure network platform for private data retrieval, and help users master their own secure data.

# Its value

FlowShield aims to build a decentralized private data retrieval security network system of web3 in the world, and help users recapture the privacy security information eroded by giants under web2, so that the current global hot zero-trust security network technology combined with web3 can better help users master their own security privacy data and give users a good experience of private data retrieval security network products.

<img  alt="image" src="https://res.cloudinary.com/malloc/image/upload/v1669810897/FlowShield/flowshield_project_desp_vjzury.png">

# How do we build it?

The design part of FlowShield project includes distributed full-nodes, network miner provider, intelligent contract, network quality checker and network client program. The details are as follows:

## Part one:[FlowShield-Fullnode](./fullnode)(Ful nodes of private data retrieval network based on DAO Tools)

Anyone can run Fullnode, which hosts the metadata of decentralized network and provides metadata networking and transaction matching platform. It integrates metadata from all providers, and providers use [libp2p-based pubsub](https://github.com/libp2p/go-libp2p) every few seconds to keep heartbeat to Fullnode to prove that they are online.

Users can find resources and nodes to build their own secure anonymous network tunnel. They only need to pay some tokens, and the provider nodes can get these tokens as rewards.

For all users' and Dao's data, we use Filecoin's web3.storage to store user data in a decentralized way.

<img width="1425" alt="image" src="https://user-images.githubusercontent.com/34047788/191491783-840a042d-4f39-4247-ae74-86e9278ebb4f.png">


<img width="1102" alt="image" src="https://user-images.githubusercontent.com/34047788/191491199-61b73816-5538-460c-b0ba-e9b662e8681d.png">


## Part two: [FlowShield-Provider](./provider)(Network Miner, a Secure Network Tunnel Provider for Decentralized Data Private Retrieval)

Our nodes are automatically networked through kademlia DHT and IPFS networks of libp2p through peer discovery and routing, and data synchronization among multiple nodes is realized through PubSub function of libp2p.

For all users and Dao data, we use web3.storage of Filecoin to store user data in a decentralized way.

<img width="1120" alt="image" src="https://user-images.githubusercontent.com/34047788/191491394-6dccc868-ed08-483b-9a74-2fcff6a243e2.png">


## Part three:[FlowShield-Contracts](./contract)(support:EVM Chains)

We provide a complete smart contract for the decentralized trusted bandwidth market. Our smart contract is deployed on the EVM network, and we provide many methods in the smart contract to ensure a safe trading process and a safe trading environment.

## Part four:[FlowShield-verifier](./verifier)(Decentralized network quality checker)
We provide the verifier component for the decentralized trusted bandwidth market. Anyone can run the network verifier, monitor the network quality of ongoing orders, and detect and punish illegal and bad network providers.

<img width="1428" alt="image" src="https://user-images.githubusercontent.com/34047788/191491491-cde176f0-f01c-4dfe-8d5f-b6f7d8964f35.png">


## Part five:[DeCA](./ca)(Decentralized PKI CA center)


Decentralize PKI CA center to provide communication authentication infrastructure for Dao point-to-point communication.

<img width="1363" alt="image" src="https://user-images.githubusercontent.com/34047788/191492613-b5b76237-38bb-468f-b15a-860f67581818.png">



## Part six:[FlowShield-Client](./client)(client for private data retrieval.)

A client user connects to a provider to establish a network security tunnel for private data retrieval.

# Disclaimers
When you use this software, you have agreed and declared that the author, maintainer and contributor of this software will not be responsible for any risks, costs or problems you encounter. If you find software defects or bugs, please submit patches to help improve!


# Thanks supports

<table>
  <tr>
    <td align="center"><a href="https://protocol.ai/"><img src="https://user-images.githubusercontent.com/34047788/188373221-4819fd05-ef2f-4e53-b784-dcfffe9c018c.png" width="100px;" alt="Protocol Labs"/><br /><sub><b>Protocol Labs</b></sub></a></td>
    <td align="center"><a href="https://filecoin.io/"><img src="https://user-images.githubusercontent.com/34047788/188373584-e245e0bb-8a3c-4773-a741-17e4023bde65.png" width="100px;" alt="Filecoin"/><br /><sub><b>Filecoin</b></sub></a></td>
     <td align="center"><a href="https://www.nervos.org/"><img src="https://user-images.githubusercontent.com/34047788/188373709-4c6caff6-be9f-497a-9bc3-88e6ae7195ac.png" width="100px;" alt="Nervos"/><br /><sub><b>Nervos</b></sub></a></td>
  </tr>
</table>

