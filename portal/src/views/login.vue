<template>
  <div class="login-page">
    <div class="mt-15">
      <v-btn x-large :href="github" rounded>
        <v-icon class="mr-5">mdi-github</v-icon>
        Sign in with Github
      </v-btn>
    </div>
    <div class="mt-15">
      <v-btn x-large rounded @click="connectWallet" id="walletBtn">
        <v-icon class="mr-5">mdi-wallet</v-icon>
        Sign in with Wallet (DID)
      </v-btn>
    </div>
<!--
    <div class="mt-15">
      <p class="mt-10" v-if="address">{{ address }}</p>
      <v-btn x-large rounded @click="connectWallet" v-else>
        <v-icon class="mr-5">mdi-wallet</v-icon>
        Connect Your Wallet
      </v-btn>
    </div>
-->
  </div>
</template>
<script>
import { getGithubIdOnCeramic } from '@/utils/ceramic'
// import store from '@/store'

export default {
  data: () => ({
    github: process.env.VUE_APP_BASE_URL + '/user/login',
    address: ''
  }),
  methods: {
    async connectWallet() {
      const profile = await getGithubIdOnCeramic()
      console.log(profile)
      if (profile && profile.githubID) {
        this.$store.commit('SET_CERAMIC', { uuid: profile.githubID, address: profile.address })
      } else {
        this.$message.error('The wallet address has not been bound to the github account')
        this.$store.commit('SET_CERAMIC', { uuid: 'false', address: '' })
      }
    }
  }
}
</script>
<style lang="scss" scoped>
.login-page {
  text-align: center;
}
</style>
