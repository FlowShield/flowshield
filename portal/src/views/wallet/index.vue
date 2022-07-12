<template>
  <div>
    <h3 class="font-weight-thin text-h3 mt-10">My Wallet</h3>
    <div class="login-page">
    <div class="mt-15">
      <div v-if="address">
        <p class="mt-10">Your wallet address is '{{ address }}'</p>
        <v-btn x-large rounded @click="changeWallet">
          <v-icon class="mr-5">mdi-wallet</v-icon>
          Change Your Wallet
        </v-btn>
      </div>
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
import { bindWallet, getWallet, changeWallet } from '../../utils/store.js'
import store from '@/store'

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
    this.getBind()
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
    async getBind() {
      const address = await getWallet(store.state.user.uuid)
      if (address !== '0x0000000000000000000000000000000000000000') {
        this.address = address
      }
    },
    async connectWallet() {
      let address = await getWallet(store.state.user.uuid)
      if (address === '0x0000000000000000000000000000000000000000') {
        await bindWallet(store.state.user.uuid)
        address = await getWallet(store.state.user.uuid)
        if (address === '0x0000000000000000000000000000000000000000') {
          this.$message.error('Bind failed')
        }
      }
    },
    async changeWallet() {
      await changeWallet(store.state.user.uuid)
      this.address = await getWallet(store.state.user.uuid)
    }
  }
}
</script>
<style lang="scss" scoped>
.login-page {
  text-align: center;
}
</style>
