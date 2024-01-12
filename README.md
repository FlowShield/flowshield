<p align="center">
<img width="800" alt="image" src="https://user-images.githubusercontent.com/34047788/204796156-c5c9d228-725b-4e14-b8ed-1e705ba19bc5.png">

</p>

# Our mission
<p align="center">

<img  alt="image" src="https://user-images.githubusercontent.com/34047788/204788311-9db29be8-ac90-4bff-9a39-6f9126ba22c7.png" width="80%">

<b> Build a global web3 decentralized private retrieval of data security network，Building Cyber  Sovereignty.</b>
</p>

In the field of new generation network security communications, the FlowShield project is leading a revolution. With the rapid development of the digital age, we are facing more and more network security threats and privacy issues. Traditional secure communications solutions are no longer able to meet the growing demands, so we need an innovative approach to protect our communications and data.

FlowShield adopts a series of cutting-edge technologies, such as blockchain, encryption algorithms and secure network architecture, to provide a strong guarantee for the new generation of network secure communication. Through decentralized design, FlowShield breaks the traditional centralized model and makes communication more secure and reliable. Each participant can jointly build and maintain this network, contributing to the security of the entire ecosystem.

FlowShield not only focuses on data encryption and transmission security, but also on privacy protection. The project is committed to protecting the personal privacy of users and preventing their sensitive information from being misused or leaked. Through an anonymized and decentralized approach, FlowShield ensures that users’ identities and communication content are protected to the greatest extent possible, allowing them to have greater autonomy and control in the digital world.

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

DeCA has advantages at every stage of the PKI life cycle. It makes the autonomous control of online identity possible, and provides a simple and more powerful SSL certificate. In use, it can help entities finally store encrypted data by degrading public key management to security decentralization.

<img width="1363" alt="image" src="https://user-images.githubusercontent.com/34047788/191492613-b5b76237-38bb-468f-b15a-860f67581818.png">

### [More documentation on DeCA](https://www.flowshield.xyz/flowshield_docs/cloudslit/deca/architecture/)



## Part six:[FlowShield-Client](./client)(client for private data retrieval.)

A client user connects to a provider to establish a network security tunnel for private data retrieval.

# Disclaimers
When you use this software, you have agreed and declared that the author, maintainer and contributor of this software will not be responsible for any risks, costs or problems you encounter. If you find software defects or bugs, please submit patches to help improve!

# Thanks supports

<table>
  <tr>
    <td align="center"><a href="https://protocol.ai/"><img src="https://user-images.githubusercontent.com/34047788/188373221-4819fd05-ef2f-4e53-b784-dcfffe9c018c.png" width="100px;" alt="Protocol Labs"/><br /><sub><b>Protocol Labs</b></sub></a></td>
    <td align="center"><a href="https://filecoin.io/"><img src="https://user-images.githubusercontent.com/34047788/188373584-e245e0bb-8a3c-4773-a741-17e4023bde65.png" width="100px;" alt="Filecoin"/><br /><sub><b>Filecoin</b></sub></a></td>
     <td align="center"><a href="https://fvm.filecoin.io/"><img src="https://user-images.githubusercontent.com/34047788/220075045-48286b37-b708-4ecf-94f5-064c55e79fa3.png" width="110px;" alt="FVM"/><br /><sub><b>FVM</b></sub></a></td>
     <td align="center"><a href="https://www.nervos.org/"><img src="https://user-images.githubusercontent.com/34047788/188373709-4c6caff6-be9f-497a-9bc3-88e6ae7195ac.png" width="100px;" alt="Nervos"/><br /><sub><b>Nervos</b></sub></a></td>
  </tr>
</table>

