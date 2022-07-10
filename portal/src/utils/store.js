import Vue from 'vue'
import * as Web3 from 'web3'
import * as contractJSON from '../../../contract/artifacts/contracts/Token.sol/Token.json'
const state = Vue.observable({ account: '', balance: 0, status: '' })

export const _connect = async() => {
  if (window.ethereum != null) {
    try {
      await window.ethereum.send('eth_requestAccounts')
      console.log(window.ethereum)
      const web3 = await new Web3(window.ethereum)
      const accounts = await web3.eth.getAccounts()
      console.log(accounts[0])
      if (accounts[0] !== undefined) {
        state.account = accounts[0]
        const balance = await web3.eth.getBalance(accounts[0])
        state.balance = balance * 0.0000000000000000001
        localStorage.setItem('connected', state.account)
      }
    } catch (error) {
      // err
    }
  } else {
    alert('Please install Metamask')
  }
}
export const _disconnect = () => {
  localStorage.removeItem('connected')
  state.account = ''
  state.balance = 0
  location.reload()
}
export const checkAndGo = async() => {
  if (window.ethereum) {
    const web3 = await new Web3(window.ethereum)
    const accounts = await web3.eth.getAccounts()
    const cache = localStorage.getItem('connected')
    if (accounts) {
      if (accounts[0] === cache) {
        state.account = accounts[0]
        const balance = await web3.eth.getBalance(accounts[0])
        state.balance = balance * 0.000000000000000001
      }
    }
  }
}
const networks = {
  dev: {
    httpProvider: 'https://ropsten.infura.io/v3/811238fc53164a35a96f841a7a89bea5'
  }
}

const contractAddress = '0xCAdF7B1f6f4E5452FAdf55820B2EA7b6Dd3C972a'

const web3Init = async() => {
  let web3 = ''
  if (typeof web3 !== 'undefined') {
    web3 = await new Web3(window.ethereum)
  } else {
    // set the provider you want from Web3.providers
    web3 = await new Web3(new Web3.providers.HttpProvider(networks.dev.httpProvider))
    console.log('null', web3)
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
  state.status = x
  console.log(x)
}

export const setStatus = async() => {
  const web3 = await web3Init()
  const accounts = await web3.eth.getAccounts()
  web3.eth.defaultAccount = accounts[0]
  const account = accounts[0]
  const abi = contractJSON.abi
  const sbs = new web3.eth.Contract(abi, contractAddress, account)
  await sbs.methods.clientOrder('test', 1000).send({ from: accounts[0] })
  getBalance()
}

export const payOrder = async(uuid, price) => {
  const web3 = await web3Init()
  const accounts = await web3.eth.getAccounts()
  web3.eth.defaultAccount = accounts[0]
  const account = accounts[0]
  const abi = contractJSON.abi
  const sbs = new web3.eth.Contract(abi, contractAddress, account)
  if (await getOrder(uuid) === true) {
    return 'Paid, payment failed'
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

export default state
