import * as Web3 from 'web3'
import * as contractJSON from '../../../contract/artifacts/contracts/CloudSlit.sol/CloudSlit.json'

export const paid = 'Paid, payment failed'

export const _connect = async() => {
  if (window.ethereum != null) {
    try {
      await window.ethereum.send('eth_requestAccounts')
      console.log(window.ethereum)
      const web3 = await new Web3(window.ethereum)
      const accounts = await web3.eth.getAccounts()
      console.log(accounts[0])
    } catch (error) {
      // err
    }
  } else {
    alert('Please install Metamask')
  }
}

export const checkAndGo = async() => {
  if (window.ethereum) {
    const web3 = await new Web3(window.ethereum)
    const accounts = await web3.eth.getAccounts()
    const cache = localStorage.getItem('connected')
    if (accounts) {
      if (accounts[0] === cache) {
        const balance = await web3.eth.getBalance(accounts[0])
      }
    }
  }
}

const networks = {
  dev: {
    httpProvider: 'https://ropsten.infura.io/v3/811238fc53164a35a96f841a7a89bea5'
  }
}

const contractAddress = '0x5c5286B0B87B075563e6308CC9D5d9D9E0901228'

const web3Init = async() => {
  let web3 = ''
  console.log(typeof web3)
  if (typeof web3 !== 'undefined') {
    web3 = await new Web3(window.ethereum)
  } else {
    // set the provider you want from Web3.providers
    web3 = await new Web3(new Web3.providers.HttpProvider(networks.dev.httpProvider))
  }
  console.log(web3)
  return web3
}

export const getBalance = async() => {
  const web3 = await web3Init()
  const accounts = await web3.eth.getAccounts()
  web3.eth.defaultAccount = accounts[0]
  const account = accounts[0]
  const abi = contractJSON.abi
  const sbs = new web3.eth.Contract(abi, contractAddress, account)
  const x = await sbs.methods.balanceOf(account).call()
  console.log(x)
}

export const getWallet = async(uid) => {
  const web3 = await web3Init()
  const accounts = await web3.eth.getAccounts()
  web3.eth.defaultAccount = accounts[0]
  const account = accounts[0]
  const abi = contractJSON.abi
  const sbs = new web3.eth.Contract(abi, contractAddress, account)
  const walletAddress = await sbs.methods.getWallet(uid).call()
  return walletAddress
}

export const bindWallet = async(uid) => {
  const web3 = await web3Init()
  const accounts = await web3.eth.getAccounts()
  web3.eth.defaultAccount = accounts[0]
  const account = accounts[0]
  const abi = contractJSON.abi
  const sbs = new web3.eth.Contract(abi, contractAddress, account)
  await sbs.methods.bindWallet(uid).send({ from: accounts[0] })
  getBalance()
}

export const changeWallet = async(uid) => {
  const web3 = await web3Init()
  const accounts = await web3.eth.getAccounts()
  web3.eth.defaultAccount = accounts[0]
  const account = accounts[0]
  const abi = contractJSON.abi
  const sbs = new web3.eth.Contract(abi, contractAddress, account)
  console.log(accounts[0])
  await sbs.methods.changeWallet(uid).send({ from: accounts[0] })
  getBalance()
}

export const payOrder = async(uuid, price) => {
  const web3 = await web3Init()
  const accounts = await web3.eth.getAccounts()
  web3.eth.defaultAccount = accounts[0]
  const account = accounts[0]
  const abi = contractJSON.abi
  const sbs = new web3.eth.Contract(abi, contractAddress, account)
  sbs.handleRevert = true
  if (await getOrder(uuid) === true) {
    return paid
  }
  try {
    await sbs.methods.clientOrder(uuid, price).send({ from: accounts[0] })
  } catch (error) {
    console.log(error)
    return 'Payment failed'
  }
  return 'ok'
}

export const getOrder = async(uuid) => {
  const web3 = await web3Init()
  const accounts = await web3.eth.getAccounts()
  web3.eth.defaultAccount = accounts[0]
  const account = accounts[0]
  const abi = contractJSON.abi
  const sbs = new web3.eth.Contract(abi, contractAddress, account)
  const x = await sbs.methods.checkOrder(uuid).call()
  console.log(uuid, x)
  return x
}
