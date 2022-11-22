<p align="center">
<img 
    src="https://user-images.githubusercontent.com/34047788/200487286-e0661f79-b54e-4fff-b5de-d8e8139e3961.png" 
    width="251" border="0" alt="CloudSlit">
</p>

https://user-images.githubusercontent.com/34047788/203346446-6bec3b69-9037-455d-92e4-f5393b3780dc.mp4


# CloudSlit-****FVM**** build notes  [(whitepaper)](https://www.cloudslit.xyz/cloudslit_whitepaper/)


- [CloudSlit-****FVM**** build notes  (whitepaper)](#cloudslit-fvm-build-notes--whitepaper)
  - [Project Description](#project-description)
  - [Project **FVM** version links](#project-fvm-version-links)
  - [How is **FVM**, Filecoin, IPFS technology used in this project](#how-is-fvm-filecoin-ipfs-technology-used-in-this-project)
    - [CloudSlit-Fullnode - build with IPFS|libp2p](#cloudslit-fullnode---build-with-ipfslibp2p)
    - [CloudSlit-Provider - build with IPFS|libp2p](#cloudslit-provider---build-with-ipfslibp2p)
    - [DeCA - build with IPFS|libp2p|**FVM**](#deca---build-with-ipfslibp2pfvm)
    - [Contracts - build with **FVM**](#contracts---build-with-fvm)
  - [**FVM**'s current work contents](#fvms-current-work-contents)
    - [Cloudslit-contracts **FVM** adaptation process record](#cloudslit-contracts-fvm-adaptation-process-record)
    - [Deploy the contract on `wallaby`](#deploy-the-contract-on-wallaby)
    - [Cloudslit-provider **FVM** adaptation process record](#cloudslit-provider-fvm-adaptation-process-record)

## Project Description

CloudSlit aims to build a decentralized web3 privacy data retrieval security network system around the world, and help users recapture privacy security information eroded by giants under web2, so that the current global hot zero-trust security network technology combined with web3 can better help users master their own security privacy data and give users a good experience of privacy data retrieval security network products.

## Project **FVM** version links
 
Link to Github repo: https://github.com/CloudSlit/cloudslit/tree/fevm-testnet

**FVM** testnet website: https://fvm-testnet.cloudslit.xyz/

## How is **FVM**, Filecoin, IPFS technology used in this project

### CloudSlit-Fullnode - build with IPFS|libp2p

All nodes are CloudSlit networks, which provide DAO tools for institutions\organizations\individuals, mainly to satisfy users’ privatization and customize their own network privatization attributes.

Anyone can run all nodes, host metadata of decentralized network, and provide metadata networking and transaction matching platform. It integrates metadata from all providers, and providers use libp2p-based pubsub every few seconds to keep their hearts beating to Fullnode to prove that they are online.

Users can find resources and nodes inside the all-node platform to build their own secure anonymous network tunnel. They only need to pay some tokens, and the network providers (miners) nodes can get these tokens as rewards.

For all users’ and Dao’s data, we store it on the decentralized network of Filecoin.

![image](https://user-images.githubusercontent.com/34047788/191491783-840a042d-4f39-4247-ae74-86e9278ebb4f.png)



### CloudSlit-Provider - build with IPFS|libp2p

Provider, as a secure network tunnel provider for decentralized private data retrieval, provides decentralized network access services for CloudSlit network users. The Provider decentralized networking through peer-to-peer discovery and routing through libp2p, and combined with the publish and subscribe function of P2P to achieve data synchronization among multiple nodes.

![image](https://user-images.githubusercontent.com/34047788/191491394-6dccc868-ed08-483b-9a74-2fcff6a243e2.png)


### DeCA - build with IPFS|libp2p|**FVM**

Decentralized PKI CA center provides communication authentication function for point-to-point communication between client and miner nodes in CloudSlit network. DeCA can perform all the key functions of X.509 PKI standard, that is, register, confirm, revoke and verify mTLS certificates.


![image](https://user-images.githubusercontent.com/34047788/191492613-b5b76237-38bb-468f-b15a-860f67581818.png)

### Contracts - build with **FVM**

CloudSlit mainly uses smart contracts to build a decentralized storage engine policy center. Our goal is to establish a private data retrieval platform that can run by itself and be managed by the public. At present, the operation carrier of smart contract mainly considers the virtual machine environment compatible with EVM. First, we choose **FVM** as our decentralized management platform. As the computing layer of the FileCoin storage ecosystem, **FVM** allows us to conduct trusted computing, provide services closer to data storage, and provide users with more reliable data computing credibility.

The main functions of smart contracts include:

1. Pledge and redemption of fullnode node and provider node
2. Matching and payment of users’ online orders
3. Withdrawal of benefits by network providers
  
## **FVM**'s current work contents

1. Adapting cloudslit fullnode, network provider, smart contracts on fevm test network.
   
2. Deploy cloudslit-FVM's exclusive website, and let cloudslit platform run on FVM's test network.
   
3. Upgrade cloudslit provider and fullnode pledge process to adapt to **FVM** test network.


### Cloudslit-contracts **FVM** adaptation process record
### Deploy the contract on `wallaby`

Our contract warehouse address is [https://github.com/CloudSlit/cloudslit/tree/fevm-testnet/contract](https://github.com/CloudSlit/cloudslit/tree/fevm-testnet/contract)

Refer to FEVM-Hardhat-Kit https://github.com/filecoin-project/FEVM-Hardhat-Kit

1. Create a .env file, fill in the private key address `PRIVATE_KEY="abcdefg"`, first execute `yarn install` to install dependencies,

Then get the f4 address `yarn hardhat get-address`, where f4address is the file currency representation of your Ethereum address.

```jsx
➜  contract git:(fevm-testnet) ✗ yarn hardhat get-address    
yarn run v1.22.17
$ /Users/chengqiang/icode/github/CloudSlit/cloudslit/contract/node_modules/.bin/hardhat get-address
You have both ethereum-waffle and @nomicfoundation/hardhat-chai-matchers installed. They don't work correctly together, so please make sure you only use one.

We recommend you migrate to @nomicfoundation/hardhat-chai-matchers. Learn how to do it here: https://hardhat.org/migrate-from-waffle
f4address =  f410fcyr4jy3t7ah2pm6v4rwc64n4kbyixjnjyirzrdq
Ethereum address: 0x1623c4E373f80fa7B3d5E46c2F71bc50708bA5A9
✨  Done in 5.66s.
```

According to the f4address address or the Ethereum address, the faucet recharge is carried out, and the deployment is carried out after the recharge is successful.

```jsx
➜  contract git:(fevm-testnet) ✗ yarn hardhat deploy
yarn run v1.22.17
$ /Users/chengqiang/icode/github/CloudSlit/cloudslit/contract/node_modules/.bin/hardhat deploy
You have both ethereum-waffle and @nomicfoundation/hardhat-chai-matchers installed. They don't work correctly together, so please make sure you only use one.

We recommend you migrate to @nomicfoundation/hardhat-chai-matchers. Learn how to do it here: https://hardhat.org/migrate-from-waffle
Compiled 2 Solidity files successfully
Wallet Ethereum Address: 0x1623c4E373f80fa7B3d5E46c2F71bc50708bA5A9
Wallet f4Address:  f410fcyr4jy3t7ah2pm6v4rwc64n4kbyixjnjyirzrdq
deploying "CloudSlitDao" (tx: 0xdfa457bbb69b7d317236d6f6665ca34bfba7d971d33da34c52b6c6ba6a054864)...: deployed at 0x7D9208Aa55833A9F5641cD2d7477733D819435eC with 70991164 gas
✨  Done in 76.95s.
```

After successful deployment, use go eth sdk and ethers.js to test

Test with ethers.js
I encountered this problem at the beginning, and the research found that the request threw an abnormal error when the contract require was judged, and the thrown error message was not displayed here.
<img width="1511" alt="image" src="https://user-images.githubusercontent.com/12872991/200886049-a3d3d3dd-874b-4e52-891f-f324aea45129.png">

Then make a call to successfully wake up MetaMask and make a signature payment, but this error will be printed out later. He will affect my call status acquisition, I don't know whether he is successful, although when he throws an error, the contract interaction has been completed.
<img width="1511" alt="image" src="https://user-images.githubusercontent.com/12872991/200897512-ec059a1e-342a-4f38-8130-bb354c8769a5.png">

### Cloudslit-provider **FVM** adaptation process record

1. Failed to stake   | ERRO[0012] runETH error : missing required field 'logsBloom' for Header
2. Occasional service unavailability | ERRO[0000] runETH error : 503 Service Unavailable: {"jsonrpc":"2.0","id":3,"error":{"code":0,"message":"fatal error calling 'eth_getBlockByNumber': panic in rpc method 'eth_getBlockByNumber': runtime error: invalid memory address or nil pointer dereference"}}

