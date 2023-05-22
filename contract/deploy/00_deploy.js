require("hardhat-deploy")
require("hardhat-deploy-ethers")


const private_key = network.config.accounts[0]
const wallet = new ethers.Wallet(private_key, ethers.provider)

module.exports = async ({ deployments }) => {
  console.log("Wallet Ethereum Address:", wallet.address)

  //deploy FlowShieldDao
  const FlowShieldDao = await ethers.getContractFactory('FlowShieldDao', wallet);
  console.log('Deploying FlowShieldDao...');
  const certificate = await FlowShieldDao.deploy();
  await certificate.deployed()
  console.log('FlowShieldDao deployed to:', certificate.address);
}