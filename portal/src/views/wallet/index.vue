<template>
  <div>
    <h3 class="font-weight-thin text-h3 mt-10">My Wallet</h3>
    <div class="login-page">
      <div class="mt-15">
        <div v-if="address">
          <template>
            <v-card
                class="mx-auto"
            >
              <v-container fluid>
                <v-row dense>
                  <v-col
                      :cols="6"
                  >
                    <v-card  dark>
                        <v-card-title >
                          Account address
<!--                          <span v-if="status == 1" style="color: #0D8DF1">-->
<!--                            &nbsp;(Awaiting verification)-->
<!--                          </span>-->
                          <span style="color: #0DF171">
                            &nbsp;(Verified)
                          </span>
                        </v-card-title>
                      <div class="text--primary">
                        {{ address }}
                      </div>
                      <v-card-actions>
                        <v-spacer></v-spacer>
                      </v-card-actions>
                      <div style="padding-bottom: 20px">
                        <v-btn x-large rounded color="primary" @click="changeWallet" :loading="changewalletLoading">
                          <v-icon>mdi-wallet</v-icon>
                          Change Your Wallet
                        </v-btn>
                      </div>
                    </v-card>
                  </v-col>
                  <v-col
                      :cols="6"
                  >
                    <v-card>
                      <v-card-title >Withdrawable CSD</v-card-title>
                      <div class="text--primary">
                        {{ withdrawCSD }}
                      </div>
                      <v-card-actions>
                        <v-spacer></v-spacer>
                      </v-card-actions>
                      <div style="padding-bottom: 20px">
                        <v-btn
                            x-large
                            rounded
                            color="teal"
                            @click="withdrawAllOrder"
                            :loading="withdrawLoading"
                        >
                          <v-icon class="mr-3">mdi-wallet</v-icon>
                          Withdraw all CSD
                        </v-btn>
                      </div>
                    </v-card>
                  </v-col>
                </v-row>
              </v-container>
            </v-card>
          </template>
        </div>
        <v-btn x-large rounded @click="bindWallet" :loading="bindwalletLoading" v-else>
          <v-icon class="mr-5">mdi-wallet</v-icon>
          Connect Your Wallet (DID)
        </v-btn>
      </div>
    </div>
  </div>

</template>
<script>
import {
  getAllOrderTokens,
  withdrawAllOrderTokens
} from '@/utils/ethers'
import store from '@/store'
// import FormDialog from './components/form-dialog'
import { updateGithubIdOnCeramic } from '@/utils/ceramic'

export default {
  components: {},
  data: () => ({
    address: '',
    withdrawCSD: 0,
    status: 0,
    color: '',
    bindwalletLoading: false,
    changewalletLoading: false,
    walletLoading: false,
    withdrawLoading: false,
    query: {
      name: '',
      page: 1,
      limit_num: 15
    },
    tableHeaders: [
      { text: 'Wallet', value: 'peer_id', width: '210px' },
      { text: 'Type', value: 'type' },
      { text: 'Loc', value: 'loc' },
      { text: 'IP', align: 'start', value: 'ip' },
      { text: 'Addr', value: 'addr' },
      { text: 'Listen port', value: 'port' },
      { text: 'Colo', value: 'colo' },
      { text: 'Gas', value: 'gas_price' },
      { text: 'Created at', value: 'CreatedAt' },
      { text: 'Updated at', value: 'UpdatedAt' }
    ],
    tableItems: [],
    total: 0
  }),
  created() {
    this.getBind()
    this.getAllOrderTokens()
  },
  methods: {
    async getBind() {
      if (store.state.ceramic.address) {
        this.address = store.state.ceramic.address
      }
    },
    async getAllOrderTokens() {
      this.withdrawCSD = await getAllOrderTokens()
    },
    async bindWallet() {
      this.bindwalletLoading = true
      const profile = await updateGithubIdOnCeramic(store.state.user.uuid)
      if (profile && profile.address) {
        this.$store.commit('SET_CERAMIC', { uuid: profile.githubID, address: profile.address })
        this.address = profile.address
      } else {
        this.$message.error('Bind Wallet error')
      }
      await this.getAllOrderTokens()
      this.bindwalletLoading = false
    },
    async changeWallet() {
      this.changewalletLoading = true
      const profile = await updateGithubIdOnCeramic(store.state.user.uuid)
      if (profile && profile.address) {
        this.$store.commit('SET_CERAMIC', { uuid: profile.githubID, address: profile.address })
        this.address = profile.address
      } else {
        this.$message.error('Bind Wallet error')
      }
      await this.getAllOrderTokens()
      this.changewalletLoading = false
    },
    async withdrawAllOrder() {
      this.withdrawLoading = true
      await withdrawAllOrderTokens()
      await this.getAllOrderTokens()
      this.withdrawLoading = false
    }
  }
}
</script>
<style lang="scss" scoped>
.login-page {
  text-align: center;
}
</style>
