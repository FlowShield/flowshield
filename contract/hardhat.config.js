require("@nomicfoundation/hardhat-toolbox");
require('hardhat-deploy');
require('hardhat-deploy-ethers');
require("./tasks")
require("dotenv").config()

const PRIVATE_KEY = process.env.PRIVATE_KEY
/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
    solidity: "0.8.17",
    networks: {
        wallaby: {
            url: "https://wallaby.node.glif.io/rpc/v0",
            accounts: [PRIVATE_KEY],
        }
    },
    defaultNetwork: "wallaby",
};