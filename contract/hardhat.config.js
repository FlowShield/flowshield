require("@nomiclabs/hardhat-waffle");

// Go to https://infura.io/ and create a new project
// Replace this with your Infura project ID
const INFURA_PROJECT_ID = "811238fc53164a35a96f841a7a89bea5";

// Replace this private key with your Ropsten account private key
// To export your private key from Metamask, open Metamask and go to Account Details > Export Private Key
// Be aware of NEVER putting real Ether into testing accounts
const ROPSTEN_PRIVATE_KEY = "a5042f010a7d7f5652097768612265014ec390ea6f2f281f362091a5c39f4900";
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