<template>
  <div>
    <h3 class="font-weight-thin text-h3 mt-10">My Wallet</h3>
    <div class="login-page">
    <div class="mt-15">
      <p class="mt-10" v-if="address">{{ address }}</p>
      <v-btn x-large rounded @click="connectWallet" v-else>
        <v-icon class="mr-5">mdi-wallet</v-icon>
        Connect Your Wallet
      </v-btn>
    </div>
    </div>
  </div>

</template>
<script>
import { fetchZeroAccessNodes } from '@/api'
import { ethers } from 'ethers'
import { getBalance, setStatus } from '../../utils/store.js'

export default {
  components: { },
  data: () => ({
    address: '',
    loading: false,
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
    this.getTableItems()
  },
  methods: {
    handleSearch() {
      this.query.page = 1
      this.getTableItems()
    },
    getTableItems() {
      this.loading = true
      fetchZeroAccessNodes(this.query).then(res => {
        this.tableItems = res.data.list || []
        this.total = res.data.paginate.total
      }).finally(() => {
        this.loading = false
      })
    },
    handleCount(v) {
      this.query.limit_num = v
      this.handleSearch()
    },
    async connectWallet() {
      const provider = new ethers.providers.Web3Provider(window.ethereum)
      // eth_requestAccounts can silent prompt
      await provider.send('wallet_requestPermissions', [{ // prompts every time
        eth_accounts: {}
      }])
      const signer = provider.getSigner()
      this.address = await signer.getAddress()
      getBalance(this.newStatus)
      setStatus()
    }
  }
}
</script>
<style lang="scss" scoped>
.login-page {
  text-align: center;
}
</style>
