const abi = require('../../../contract/artifacts/contracts/CloudSlit.sol/CloudSlit.json').abi
const ethers = require('ethers')
export const OrderPaid = 'Paid, payment failed'
export const BalanceNotEnough = 'Your balance is insufficient'

const providerInit = async() => {
  let provider = ''
  if (typeof window.ethereum !== 'undefined') {
    provider = new ethers.providers.Web3Provider(window.ethereum)
  } else {
    // set the provider you want from Web3.providers
    alert('Please install MetaMask')
  }
  return provider
}

const contractAddress = '0xa9eC0d7FF7975D9930030DDf3d0c9260B1A766cF'

export const getBalance = async() => {
  const provider = await providerInit()
  // eth_requestAccounts can silent prompt
  await provider.send('eth_requestAccounts', [])
  const signer = provider.getSigner()
  const address = await signer.getAddress()
  const contract = new ethers.Contract(contractAddress, abi, provider)
  const balance = await contract.balanceOf(address)
  return balance
}

export const getWallet = async(uid) => {
  const provider = await providerInit()
  const contract = new ethers.Contract(contractAddress, abi, provider)
  const walletAddress = await contract.getWallet(uid)
  return walletAddress
}

export const bindWallet = async(uid) => {
  const provider = await providerInit()
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  try {
    const transaction = await contract.bindWallet(uid)
    await transaction.wait()
  } catch (error) {
    console.log(error.message.message)
    return error.message.message
  }
}

export const changeWallet = async(uid, newwallet) => {
  const provider = await providerInit()
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  try {
    const transaction = await contract.changeWallet(uid, newwallet)
    await transaction.wait()
  } catch (error) {
    console.log(error)
    return error.message.message
  }
}

export const payOrder = async(name, duration, uuid, price, wallet) => {
  const provider = await providerInit()
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  if (await checkOrder(uuid) === true) {
    return OrderPaid
  }
  if (await getBalance() < price) {
    return BalanceNotEnough
  }
  try {
    const transaction = await contract.clientOrder(name, duration, uuid, price, wallet)
    await transaction.wait()
  } catch (error) {
    console.log(error.error.message)
    return 'Payment failed'
  }
  return 'ok'
}

export const checkOrder = async(uuid) => {
  const provider = await providerInit()
  const contract = new ethers.Contract(contractAddress, abi, provider)
  const orderStatus = await contract.checkOrder(uuid)
  return orderStatus
}

export const getOrdersInfo = async(uuid) => {
  const provider = await providerInit()
  const contract = new ethers.Contract(contractAddress, abi, provider)
  const orderStatus = await contract.getOrdersInfo(uuid)
  return orderStatus
}

export const getAllOrderTokens = async() => {
  const provider = await providerInit()
  await provider.send('eth_requestAccounts', [])
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  const withdrawCSD = await contract.getAllOrderTokens()
  console.log(withdrawCSD)
  return withdrawCSD
}

export const withdrawAllOrderTokens = async() => {
  const provider = await providerInit()
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  try {
    const transaction = await contract.withdrawAllOrderTokens()
    await transaction.wait()
  } catch (error) {
    console.log(error.toString())
    console.log(error.message.message)
  }
}

export const withdrawOrderTokens = async(order_id) => {
  const provider = await providerInit()
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  try {
    const transaction = await contract.withdrawOrderTokens(order_id)
    await transaction.wait()
    console.log(transaction)
  } catch (error) {
    console.log(error.toString())
    console.log(error.message)
  }
}

export const stake = async(_type) => {
  const provider = await providerInit()
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  try {
    const transaction = await contract.stake(_type)
    await transaction.wait()
  } catch (error) {
    console.log(error.toString())
    console.log(error.message.message)
  }
}
