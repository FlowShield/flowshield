require('dotenv').config();
require("@nomiclabs/hardhat-waffle");
require("@openzeppelin/hardhat-upgrades");
// Go to https://infura.io/ and create a new project
// Replace this with your Infura project ID
const INFURA_PROJECT_ID = "811238fc53164a35a96f841a7a89bea5";

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
    solidity: "0.8.15",
    networks: {
        ropsten: {
            url: `https://ropsten.infura.io/v3/${INFURA_PROJECT_ID}`,
            accounts: [process.env.PRIVATE_KEY]
        },
        'godwoken-testnet': {
            url: `https://godwoken-testnet-v1.ckbapp.dev`,
            accounts: [process.env.PRIVATE_KEY]
        },
        matic: {
            url: "https://rpc-mumbai.maticvigil.com",
            accounts: [process.env.PRIVATE_KEY]
        }
    },
    etherscan: {
        apiKey: process.env.POLYGONSCAN_API_KEY
    },
    defaultNetwork: "matic",
};