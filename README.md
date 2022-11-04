# Moralis-Filecoin-hackthon-CloudSlit

- [Moralis-Filecoin-hackthon-CloudSlit](#moralis-filecoin-hackthon-cloudslit)
- [inspiration](#inspiration)
- [Its value](#its-value)
- [How do we build it?](#how-do-we-build-it)
  - [Part one:CloudSlit-Fullnode(Ful nodes of private data retrieval network based on DAO Tools)](#part-onecloudslit-fullnodeful-nodes-of-private-data-retrieval-network-based-on-dao-tools)
- [](#)
- [inspiration](#inspiration-1)
- [Its value](#its-value-1)
- [How do we build it?](#how-do-we-build-it-1)
  - [Part one:CloudSlit-Fullnode(Ful nodes of private data retrieval network based on DAO Tools)](#part-onecloudslit-fullnodeful-nodes-of-private-data-retrieval-network-based-on-dao-tools-1)
  - [Part two: CloudSlit-Provider(Network Miner, a Secure Network Tunnel Provider for Decentralized Data Private Retrieval)](#part-two-cloudslit-providernetwork-miner-a-secure-network-tunnel-provider-for-decentralized-data-private-retrieval)
  - [Part three:CloudSlit-Contracts(support:EVM Chains)](#part-threecloudslit-contractssupportevm-chains)
  - [Part four:CloudSlit-verifier(Decentralized network quality checker)](#part-fourcloudslit-verifierdecentralized-network-quality-checker)
  - [Part five:DeCA(Decentralized PKI CA center)](#part-fivedecadecentralized-pki-ca-center)
  - [Part six:CloudSlit-Client(client for private data retrieval.)](#part-sixcloudslit-clientclient-for-private-data-retrieval)
  - [Part two: CloudSlit-Provider(Network Miner, a Secure Network Tunnel Provider for Decentralized Data Private Retrieval)](#part-two-cloudslit-providernetwork-miner-a-secure-network-tunnel-provider-for-decentralized-data-private-retrieval-1)
  - [Part three:CloudSlit-Contracts(support:polygon)](#part-threecloudslit-contractssupportpolygon)
  - [Part four:CloudSlit-verifier(Decentralized network quality checker)](#part-fourcloudslit-verifierdecentralized-network-quality-checker-1)
  - [Part five:DeCA(Decentralized PKI CA center)](#part-fivedecadecentralized-pki-ca-center-1)
  - [Part six:CloudSlit-Client(client for private data retrieval.)](#part-sixcloudslit-clientclient-for-private-data-retrieval-1)
- [The challenges we encountered](#the-challenges-we-encountered)
- [What we are proud of](#what-we-are-proud-of)


# inspiration

At present, the options available for interactive (low-latency) communication with privacy guarantee are very limited, and the solutions developed so far all focus on the traditional web model of single source data publisher, and it has defects in delay and threat models.

CloudSlit uses blockchain, web3 and secure network technology of private data retrieval to enhance and improve network security/privacy protection of users' privatization.

In order to protect the public's network security under web2, a very popular zero-trust security architecture has emerged. Our team has been working on open source products with zero trust security, but we found that although many zero trust network security companies provide zero trust security platforms, they monopolize users' network access nodes and centrally store users' core security profiles. Therefore, we are considering whether we can use web3 technology to realize a secure network for private data retrieval. We designed CloudSlit project to provide users with a decentralized secure network platform for private data retrieval, and help users master their own secure data.

# Its value

CloudSlit aims to build a decentralized private data retrieval security network system of web3 in the world, and help users recapture the privacy security information eroded by giants under web2, so that the current global hot zero-trust security network technology combined with web3 can better help users master their own security privacy data and give users a good experience of private data retrieval security network products.

<img width="1423" alt="image" src="https://user-images.githubusercontent.com/34047788/191491619-91a8030b-30f9-426d-af3a-de93739af1ad.png">

# How do we build it?

The design part of CloudSlit project includes distributed full-nodes, network miner provider, intelligent contract, network quality checker and network client program. The details are as follows:

## Part one:[CloudSlit-Fullnode](./fullnode)(Ful nodes of private data retrieval network based on DAO Tools)

Anyone can run Fullnode, which hosts the metadata of decentralized network and provides metadata networking and transaction matching platform. It integrates metadata from all providers, and providers use [libp2p-based pubsub](https://github.com/libp2p/go-libp2p) every few seconds to keep heartbeat to Fullnode to prove that they are online.

Users can find resources and nodes to build their own secure anonymous network tunnel. They only need to pay some tokens, and the provider nodes can get these tokens as rewards.

For all users' and Dao's data, we use Filecoin's web3.storage to store user data in a decentralized way.

<img width="1425" alt="image" src="https://user-images.githubusercontent.com/34047788/191491783-840a042d-4f39-4247-ae74-86e9278ebb4f.png">


<img width="1102" alt="image" src="https://user-images.githubusercontent.com/34047788/191491199-61b73816-5538-460c-b0ba-e9b662e8681d.png">
=======
  - [Part three:CloudSlit-Contracts(support:Nervos)](#part-threecloudslit-contractssupportnervos)
  - [Part four:CloudSlit-verifier(Decentralized network quality checker)](#part-fourcloudslit-verifierdecentralized-network-quality-checker)
  - [Part five:DeCA(Decentralized PKI CA center)](#part-fivedecadecentralized-pki-ca-center)
  - [Part six:CloudSlit-Client(client for private data retrieval.)](#part-sixcloudslit-clientclient-for-private-data-retrieval)
- [Disclaimers](#disclaimers)
- [Thanks supports](#thanks-supports)

# inspiration

At present, the options available for interactive (low-latency) communication with privacy guarantee are very limited, and the solutions developed so far all focus on the traditional web model of single source data publisher, and it has defects in delay and threat models.

CloudSlit uses blockchain, web3 and secure network technology of private data retrieval to enhance and improve network security/privacy protection of users' privatization.

In order to protect the public's network security under web2, a very popular zero-trust security architecture has emerged. Our team has been working on open source products with zero trust security, but we found that although many zero trust network security companies provide zero trust security platforms, they monopolize users' network access nodes and centrally store users' core security profiles. Therefore, we are considering whether we can use web3 technology to realize a secure network for private data retrieval. We designed CloudSlit project to provide users with a decentralized secure network platform for private data retrieval, and help users master their own secure data.

# Its value

CloudSlit aims to build a decentralized private data retrieval security network system of web3 in the world, and help users recapture the privacy security information eroded by giants under web2, so that the current global hot zero-trust security network technology combined with web3 can better help users master their own security privacy data and give users a good experience of private data retrieval security network products.

<img  alt="image" src="https://user-images.githubusercontent.com/34047788/191542695-3d2d6288-8a7e-4f60-b582-954942f4a907.png">

# How do we build it?

The design part of CloudSlit project includes distributed full-nodes, network miner provider, intelligent contract, network quality checker and network client program. The details are as follows:

## Part one:[CloudSlit-Fullnode](./fullnode)(Ful nodes of private data retrieval network based on DAO Tools)

Anyone can run Fullnode, which hosts the metadata of decentralized network and provides metadata networking and transaction matching platform. It integrates metadata from all providers, and providers use [libp2p-based pubsub](https://github.com/libp2p/go-libp2p) every few seconds to keep heartbeat to Fullnode to prove that they are online.

Users can find resources and nodes to build their own secure anonymous network tunnel. They only need to pay some tokens, and the provider nodes can get these tokens as rewards.

For all users' and Dao's data, we use Filecoin's web3.storage to store user data in a decentralized way.

<img width="1425" alt="image" src="https://user-images.githubusercontent.com/34047788/191491783-840a042d-4f39-4247-ae74-86e9278ebb4f.png">


<img width="1102" alt="image" src="https://user-images.githubusercontent.com/34047788/191491199-61b73816-5538-460c-b0ba-e9b662e8681d.png">


## Part two: [CloudSlit-Provider](./provider)(Network Miner, a Secure Network Tunnel Provider for Decentralized Data Private Retrieval)

Our nodes are automatically networked through kademlia DHT and IPFS networks of libp2p through peer discovery and routing, and data synchronization among multiple nodes is realized through PubSub function of libp2p.

For all users and Dao data, we use web3.storage of Filecoin to store user data in a decentralized way.

<img width="1120" alt="image" src="https://user-images.githubusercontent.com/34047788/191491394-6dccc868-ed08-483b-9a74-2fcff6a243e2.png">


## Part three:[CloudSlit-Contracts](./contract)(support:EVM Chains)

We provide a complete smart contract for the decentralized trusted bandwidth market. Our smart contract is deployed on the EVM network, and we provide many methods in the smart contract to ensure a safe trading process and a safe trading environment.

## Part four:[CloudSlit-verifier](./verifier)(Decentralized network quality checker)
We provide the verifier component for the decentralized trusted bandwidth market. Anyone can run the network verifier, monitor the network quality of ongoing orders, and detect and punish illegal and bad network providers.

<img width="1428" alt="image" src="https://user-images.githubusercontent.com/34047788/191491491-cde176f0-f01c-4dfe-8d5f-b6f7d8964f35.png">


## Part five:[DeCA](./ca)(Decentralized PKI CA center)


Decentralize PKI CA center to provide communication authentication infrastructure for Dao point-to-point communication.

<img width="1363" alt="image" src="https://user-images.githubusercontent.com/34047788/191492613-b5b76237-38bb-468f-b15a-860f67581818.png">



## Part six:[CloudSlit-Client](./client)(client for private data retrieval.)

A client user connects to a provider to establish a network security tunnel for private data retrieval.


## Part two: [CloudSlit-Provider](./provider)(Network Miner, a Secure Network Tunnel Provider for Decentralized Data Private Retrieval)

Our nodes are automatically networked through kademlia DHT and IPFS networks of libp2p through peer discovery and routing, and data synchronization among multiple nodes is realized through PubSub function of libp2p.

For all users and Dao data, we use web3.storage of Filecoin to store user data in a decentralized way.

<img width="1120" alt="image" src="https://user-images.githubusercontent.com/34047788/191491394-6dccc868-ed08-483b-9a74-2fcff6a243e2.png">


## Part three:[CloudSlit-Contracts](./contract)(support:polygon)

We provide a complete smart contract for the decentralized trusted bandwidth market. Our smart contract is deployed on the polygon test network, and we provide many methods in the smart contract to ensure a safe trading process and a safe trading environment.

Test network address: https://rpc-mumbai.maticvigil.com

Address of the contract: 0x9672f063ccba1e4ac40d31f4c00fdc9de491ab59

## Part four:[CloudSlit-verifier](./verifier)(Decentralized network quality checker)
We provide the verifier component for the decentralized trusted bandwidth market. Anyone can run the network verifier, monitor the network quality of ongoing orders, and detect and punish illegal and bad network providers.

<img width="1428" alt="image" src="https://user-images.githubusercontent.com/34047788/191491491-cde176f0-f01c-4dfe-8d5f-b6f7d8964f35.png">


## Part five:[DeCA](./ca)(Decentralized PKI CA center)


Decentralize PKI CA center to provide communication authentication infrastructure for Dao point-to-point communication.

<img width="1363" alt="image" src="https://user-images.githubusercontent.com/34047788/191492613-b5b76237-38bb-468f-b15a-860f67581818.png">



## Part six:[CloudSlit-Client](./client)(client for private data retrieval.)

A client user connects to a provider to establish a network security tunnel for private data retrieval.


# The challenges we encountered

1. Establish a stable decentralized network.

2. Construct a decentralized pubsub signaling mechanism to meet the coordination of network tunnel actions between network miners and all nodes.

3. Construct decentralized storage, encryption and decryption of user security rule data.

# What we are proud of

We have been able to satisfy users who choose decentralized network tunnel miners, create smart contract orders, pay orders, automatically build a two-way network tunnel for private retrieval of data, and enjoy the full-stack decentralized network security tunnel experience for private retrieval of data, which is our greatest joy.
