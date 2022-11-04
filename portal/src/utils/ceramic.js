// import all dependencies:
import { CeramicClient } from '@ceramicnetwork/http-client'
import { EthereumAuthProvider } from '@ceramicnetwork/blockchain-utils-linking'
import { DIDDataStore } from '@glazed/did-datastore'
import { DIDSession } from '@glazed/did-session'

const initErr = 'init err'

let walletAddress = ''

export const initCeramic = async() => {
  try {
    if (window.ethereum == null) {
      // set the provider you want from Web3.providers
      alert('Please install MetaMask')
      return initErr
    }
    // create a new CeramicClient instance:
    const ceramic = new CeramicClient('https://ceramic-clay.3boxlabs.com')

    // reference the data models this application will use:
    const aliases = {
      schemas: {
        basicProfile: 'ceramic://k3y52l7qbv1fryjn62sggjh1lpn11c56qfofzmty190d62hwk1cal1c7qc5he54ow'

      },
      definitions: {
        BasicProfile: 'kjzl6cwe1jw145cjbeko9kil8g9bxszjhyde21ob8epxuxkaon1izyqsu8wgcic'
      },
      tiles: {}
    }

    const accounts = await window.ethereum.request({
      method: 'eth_requestAccounts'
    })
    walletAddress = accounts[0]
    const authProvider = new EthereumAuthProvider(window.ethereum, accounts[0])
    const session = new DIDSession({ authProvider })

    ceramic.did = await session.authorize()

    // configure the datastore to use the ceramic instance and data models referenced above:
    return new DIDDataStore({ ceramic, model: aliases })
  } catch (error) {
    console.error(error)
    return initErr
  }
}

export const getGithubIdOnCeramic = async() => {
  try {
    const datastore = await initCeramic()
    if (datastore === initErr) {
      return
    }
    // use the DIDDatastore to get profile data from Ceramic
    return await datastore.get('BasicProfile')
  } catch (error) {
    console.error(error)
  }
}

export const updateGithubIdOnCeramic = async(uuid) => {
  try {
    const datastore = await initCeramic()
    if (datastore === initErr) {
      return
    }

    // upload
    const githubID = uuid
    const address = walletAddress

    // object needs to conform to the datamodel
    // name -> exists
    // hair-color -> DOES NOT EXIST
    const updatedProfile = {
      githubID,
      address
    }
    // use the DIDDatastore to merge profile data to Ceramic
    await datastore.merge('BasicProfile', updatedProfile)

    // use the DIDDatastore to get profile data from Ceramic
    const profile = await datastore.get('BasicProfile')
    console.log(profile)
    return profile
  } catch (error) {
    console.error(error)
  }
}
export const updateInfoOnCeramic = async() => {
  await updateGithubIdOnCeramic()
}
