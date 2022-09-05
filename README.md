# CloudSlit

#### Build a global web3 decentralized zero-trust security network，Building Cyber Sovereignty.

![image](https://user-images.githubusercontent.com/34047788/182298515-98238940-3d8a-4ce8-ab7b-c0a27666cf0a.png)

CloudSlit aims to build a global web3 decentralized zero trust security network system to help users regain the private security information eroded by giants under web2, so that the current global hot zero trust security network technology can better help users master their own security privacy data in combination with web3, and give users a good zero trust security network products and platform experience.

## Project composition：

* [CloudSlit-Fullnode](https://github.com/CloudSlit/cloudslit/tree/main/fullnode)
  (Thank Filecoin and IPFS for their support)
* [CloudSlit-Provider](https://github.com/CloudSlit/cloudslit/tree/main/provider)
  (Thank Filecoin and IPFS for their support)
* [CloudSlit-Contracts](https://github.com/CloudSlit/cloudslit/tree/main/contract)
  (Thank Nervos for their support)
* [CloudSlit-Client](https://github.com/CloudSlit/cloudslit/tree/main/client)

# Project Explanation




https://user-images.githubusercontent.com/34047788/179491684-f565275f-0e5e-44ae-b0eb-cb9a1ae11f76.mp4



# Project Design
## [CloudSlit-Fullnode](https://github.com/CloudSlit/cloudslit/tree/main/fullnode)
Anyone can run a full node, which hosts the metadata of the decentralized network, and provides a metadata networking and transaction matching platform.

For all users and Dao's data, we use Filecoin's web3.storage decentralized storage of user data.

<img width="997" alt="image" src="https://user-images.githubusercontent.com/52234994/179184171-f881f3ee-e7ca-45ad-94e1-813b9964e524.png">

## [CloudSlit-Provider](https://github.com/CloudSlit/cloudslit/tree/main/provider)
Our nodes realize automatic networking through peer discovery and routing through libp2p kademlia DHT and IPFs networks, and realize data synchronization between multiple nodes through libp2p's PubSub function.
<img width="989" alt="image" src="https://user-images.githubusercontent.com/52234994/179186444-81e0f4de-a2c1-4607-bf66-275d20c2fe0c.png">

## [CloudSlit-Contract](https://github.com/CloudSlit/cloudslit/tree/main/contract)
We use nervos to deploy smart contracts.Provide a safe trading process and a safe trading environment.

## [CloudSlit-Client](https://github.com/CloudSlit/cloudslit/tree/main/client)
The client software user connects to the provider to establish a zero trust network security tunnel.
<img width="985" alt="image" src="https://user-images.githubusercontent.com/52234994/179190148-ebd19f1d-90f0-4377-a57d-7c4942d5e0b3.png">

# Disclaimers
When you use this software, you have agreed and declared that the author, maintainer and contributor of this software will not be responsible for any risks, costs or problems you encounter. If you find software defects or bugs, please submit patches to help improve!

# Thanks for the support

<table>
  <tr>
    <td align="center"><a href="https://protocol.ai/"><img src="https://user-images.githubusercontent.com/34047788/188373221-4819fd05-ef2f-4e53-b784-dcfffe9c018c.png" width="100px;" alt=""/><br /><sub><b>Protocol Labs</b></sub></a></td>
    <td align="center"><a href="https://filecoin.io/"><img src="https://user-images.githubusercontent.com/34047788/188373584-e245e0bb-8a3c-4773-a741-17e4023bde65.png" width="100px;" alt=""/><br /><sub><b>FileCoin</b></sub></a></td>
     <td align="center"><a href="https://www.nervos.org/"><img src="https://user-images.githubusercontent.com/34047788/188373709-4c6caff6-be9f-497a-9bc3-88e6ae7195ac.png" width="100px;" alt=""/><br /><sub><b>Nervos</b></sub></a></td>
  </tr>
</table>
