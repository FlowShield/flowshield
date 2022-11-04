const abi = require('../../../contract/artifacts/contracts/CloudSlitDao.sol/CloudSlitDao.json').abi
const ethers = require('ethers')
export const OrderPaid = 'Paid, payment failed'
export const BalanceNotEnough = 'Your balance is insufficient'
export const PaySuccess = 'Payment successful'

export const Error = 'Error parsing failed'

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

const handleErrorMsg = async(error) => {
  if (error.data !== undefined && error.data.message !== undefined) {
    return error.data.message
  } else if (error.message !== undefined) {
    return error.message
  } else if (error.error !== undefined && error.error.message !== undefined) {
    return error.error.message
  } else {
    return Error
  }
}

const contractAddress = '0x9672F063Ccba1e4aC40d31f4c00fdC9dE491aB59'

export const getBalance = async() => {
  const provider = await providerInit()
  // eth_requestAccounts can silent prompt
  await provider.send('eth_requestAccounts', [])
  const signer = provider.getSigner()
  const address = await signer.getAddress()
  const contract = new ethers.Contract(contractAddress, abi, provider)
  const balance = ethers.utils.formatUnits(await contract.balanceOf(address))
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
    return handleErrorMsg(error)
  }
}
//
// export const changeWallet = async(uid, newwallet) => {
//   const provider = await providerInit()
//   const signer = provider.getSigner()
//   const contract = new ethers.Contract(contractAddress, abi, signer)
//   try {
//     const transaction = await contract.changeWallet(uid, newwallet)
//     await transaction.wait()
//   } catch (error) {
//     return handleErrorMsg(error)
//   }
// }
//
// export const verifyWallet = async(uid) => {
//   const provider = await providerInit()
//   const signer = provider.getSigner()
//   const contract = new ethers.Contract(contractAddress, abi, signer)
//   try {
//     const transaction = await contract.verifyWallet(uid)
//     await transaction.wait()
//   } catch (error) {
//     return handleErrorMsg(error)
//   }
// }

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
    const transaction = await contract.clientOrder(name, duration, uuid, ethers.utils.parseUnits(price.toString()), wallet)
    await transaction.wait()
  } catch (error) {
    return handleErrorMsg(error)
  }
  return PaySuccess
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
  const withdrawCSD = ethers.utils.formatUnits(await contract.getAllOrderTokens())
  return withdrawCSD
}

export const withdrawAllOrderTokens = async() => {
  const provider = await providerInit()
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  try {
    const transaction = await contract.withdrawAllOrderTokens()
    await transaction.wait()
    return 'Withdraw success'
  } catch (error) {
    return handleErrorMsg(error)
  }
}

export const withdrawOrderTokens = async(order_id) => {
  const provider = await providerInit()
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  try {
    const transaction = await contract.withdrawOrderTokens(order_id)
    await transaction.wait()
    return 'Withdraw success'
  } catch (error) {
    return handleErrorMsg(error)
  }
}

export const stake = async(_type) => {
  const provider = await providerInit()
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  try {
    const transaction = await contract.stake(_type)
    await transaction.wait()
    return 'Stake success'
  } catch (error) {
    return handleErrorMsg(error)
  }
}

export const isDeposit = async(_type) => {
  const provider = await providerInit()
  await provider.send('eth_requestAccounts', [])
  const signer = provider.getSigner()
  const contract = new ethers.Contract(contractAddress, abi, signer)
  const stakeStatus = await contract.isDeposit(_type)
  return stakeStatus
}
