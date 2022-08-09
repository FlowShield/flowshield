require("@nomiclabs/hardhat-waffle");
require("@openzeppelin/hardhat-upgrades");
// Go to https://infura.io/ and create a new project
// Replace this with your Infura project ID
const INFURA_PROJECT_ID = "Your project ID";

// Replace this private key with your Ropsten account private key
// To export your private key from Metamask, open Metamask and go to Account Details > Export Private Key
// Be aware of NEVER putting real Ether into testing accounts
const ROPSTEN_PRIVATE_KEY = "Your private_key";

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: "0.8.15",
  networks: {
    ropsten: {
      url: `https://ropsten.infura.io/v3/${INFURA_PROJECT_ID}`,
      accounts: [`0x${ROPSTEN_PRIVATE_KEY}`]
    },
    'godwoken-testnet': {
      url: `https://godwoken-testnet-v1.ckbapp.dev`,
      accounts: [`0x${ROPSTEN_PRIVATE_KEY}`]
    }
  },
};