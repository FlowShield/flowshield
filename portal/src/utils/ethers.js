const abi = require('../../../contract/artifacts/contracts/CloudSlit.sol/CloudSlit.json').abi
const ethers = require("ethers");



const provider = new ethers.providers.Web3Provider(window.ethereum)
// eth_requestAccounts can silent prompt
await provider.send('wallet_requestPermissions', [{ // prompts every time
  eth_accounts: {}
}])
const signer = provider.getSigner()
this.address = await signer.getAddress()
getBalance(this.newStatus)
setStatus()
