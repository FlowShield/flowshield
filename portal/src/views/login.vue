<template>
  <div class="login-page">
    <div class="mt-15">
      <v-btn x-large :href="github" rounded>
        <v-icon class="mr-5">mdi-github</v-icon>
        Sign in with Github
      </v-btn>
    </div>

    <div class="mt-15">
      <p class="mt-10" v-if="address">{{ address }}</p>
      <v-btn x-large rounded @click="connectWallet" v-else>
        <v-icon class="mr-5">mdi-wallet</v-icon>
        Connect Your Wallet
      </v-btn>
    </div>
  </div>
</template>
<script>
import { ethers } from 'ethers'

export default {
  data: () => ({
    github: process.env.VUE_APP_BASE_URL + '/user/login/github',
    address: ''
  }),
  methods: {
    async connectWallet() {
      const provider = new ethers.providers.Web3Provider(window.ethereum)
      // eth_requestAccounts can silent prompt
      await provider.send('wallet_requestPermissions', [{ // prompts every time
        eth_accounts: {}
      }])
      const signer = provider.getSigner()
      this.address = await signer.getAddress()
    }
  }
}
</script>
<style lang="scss" scoped>
.login-page {
  text-align: center;
}
</style>
