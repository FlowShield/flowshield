# Moralis-Filecoin-hackthon-CloudSlit

- [Moralis-Filecoin-hackthon-CloudSlit](#moralis-filecoin-hackthon-cloudslit)
- [inspiration](#inspiration)
- [Its value](#its-value)
- [Demo video](#demo-video)
- [How do we build it?](#how-do-we-build-it)
  - [part one:CloudSlit-Fullnode(Ful nodes of private data retrieval network based on DAO Tools)](#part-onecloudslit-fullnodeful-nodes-of-private-data-retrieval-network-based-on-dao-tools)
  - [Part two: CloudSlit-Provider(Network Miner, a Secure Network Tunnel Provider for Decentralized Data Private Retrieval)](#part-two-cloudslit-providernetwork-miner-a-secure-network-tunnel-provider-for-decentralized-data-private-retrieval)
  - [part threee:CloudSlit-Contracts(support:polygon)](#part-threeecloudslit-contractssupportpolygon)
  - [part four:CloudSlit-verifier(Decentralized network quality checker)](#part-fourcloudslit-verifierdecentralized-network-quality-checker)
  - [part five:CloudSlit-Client(client for private data retrieval.)](#part-fivecloudslit-clientclient-for-private-data-retrieval)
- [The challenges we encountered](#the-challenges-we-encountered)
- [What we are proud of](#what-we-are-proud-of)


# inspiration

At present, the options available for interactive (low-latency) communication with privacy guarantee are very limited, and the solutions developed so far all focus on the traditional web model of single source data publisher, and it has defects in delay and threat models.

CloudSlit uses blockchain, web3 and secure network technology of private data retrieval to enhance and improve network security/privacy protection of users' privatization.

In order to protect the public's network security under web2, a very popular zero-trust security architecture has emerged. Our team has been working on open source products with zero trust security, but we found that although many zero trust network security companies provide zero trust security platforms, they monopolize users' network access nodes and centrally store users' core security profiles. Therefore, we are considering whether we can use web3 technology to realize a secure network for private data retrieval. We designed CloudSlit project to provide users with a decentralized secure network platform for private data retrieval, and help users master their own secure data.

# Its value

CloudSlit aims to build a decentralized private data retrieval security network system of web3 in the world, and help users recapture the privacy security information eroded by giants under web2, so that the current global hot zero-trust security network technology combined with web3 can better help users master their own security privacy data and give users a good experience of private data retrieval security network products.

<img width="1176" alt="image" src="https://user-images.githubusercontent.com/34047788/190673496-5bbb9c23-8910-4192-9911-080230934b47.png">

# Demo video



# How do we build it?

The design part of CloudSlit project includes distributed full-nodes, network miner provider, intelligent contract, network quality checker and network client program. The details are as follows:

## part one:[CloudSlit-Fullnode](https://github.com/wanxiang-blockchain/2022-WX-Blockchain-Fall-Hackathon-CloudSlit/tree/main/codes/fullnode)(Ful nodes of private data retrieval network based on DAO Tools)

Anyone can run Fullnode, which hosts the metadata of decentralized network and provides metadata networking and transaction matching platform. It integrates metadata from all providers, and providers use [libp2p-based pubsub](https://github.com/libp2p/go-libp2p) every few seconds to keep heartbeat to Fullnode to prove that they are online.

Users can find resources and nodes to build their own secure anonymous network tunnel. They only need to pay some tokens, and the provider nodes can get these tokens as rewards.

For all users' and Dao's data, we use Filecoin's web3.storage to store user data in a decentralized way.

<img width="929" alt="image" src="https://user-images.githubusercontent.com/34047788/190649560-9bb3d443-3a3e-4747-8805-931c79db55b0.png">


## Part two: [CloudSlit-Provider](https://github.com/wanxiang-blockchain/2022-WX-Blockchain-Fall-Hackathon-CloudSlit/tree/main/codes/provider)(Network Miner, a Secure Network Tunnel Provider for Decentralized Data Private Retrieval)

Our nodes are automatically networked through kademlia DHT and IPFS networks of libp2p through peer discovery and routing, and data synchronization among multiple nodes is realized through PubSub function of libp2p.

For all users and Dao data, we use web3.storage of Filecoin to store user data in a decentralized way.

<img width="952" alt="image" src="https://user-images.githubusercontent.com/34047788/190649367-aea3230f-6c87-4b74-b692-e328ed33a78f.png">


## part threee:[CloudSlit-Contracts](https://github.com/wanxiang-blockchain/2022-WX-Blockchain-Fall-Hackathon-CloudSlit/tree/main/codes/contract)(support:polygon)

We provide a complete smart contract for the decentralized trusted bandwidth market. Our smart contract is deployed on the polygon test network, and we provide many methods in the smart contract to ensure a safe trading process and a safe trading environment.

Test network address: https://rpc-mumbai.maticvigil.com

Address of the contract: 0x9672f063ccba1e4ac40d31f4c00fdc9de491ab59

## part four:[CloudSlit-verifier](https://github.com/wanxiang-blockchain/2022-WX-Blockchain-Fall-Hackathon-CloudSlit/tree/main/codes/verifier)(Decentralized network quality checker)
We provide the verifier component for the decentralized trusted bandwidth market. Anyone can run the network verifier, monitor the network quality of ongoing orders, and detect and punish illegal and bad network providers.

<img width="935" alt="image" src="https://user-images.githubusercontent.com/34047788/190649698-9a48853b-aae4-4951-9ffe-f18ae23e8fc9.png">



## part five:[CloudSlit-Client](https://github.com/wanxiang-blockchain/2022-WX-Blockchain-Fall-Hackathon-CloudSlit/tree/main/codes/client)(client for private data retrieval.)

A client user connects to a provider to establish a network security tunnel for private data retrieval.

<img width="932" alt="image" src="https://user-images.githubusercontent.com/34047788/190649859-ee288f05-3581-4323-b672-de9546f8758d.png">



# The challenges we encountered

1. Establish a stable decentralized network.

2. Construct a decentralized pubsub signaling mechanism to meet the coordination of network tunnel actions between network miners and all nodes.

3. Construct decentralized storage, encryption and decryption of user security rule data.

# What we are proud of

We have been able to satisfy users who choose decentralized network tunnel miners, create smart contract orders, pay orders, automatically build a two-way network tunnel for private retrieval of data, and enjoy the full-stack decentralized network security tunnel experience for private retrieval of data, which is our greatest joy.