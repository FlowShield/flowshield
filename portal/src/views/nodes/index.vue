<template>
  <div>
    <h3 class="font-weight-thin text-h3 mt-10">{{ total }} Nodes</h3>
    <v-data-table
        :headers="tableHeaders"
        :items="tableItems"
        :loading="loading"
        :page="query.page"
        :items-per-page="query.limit_num"
        sort-by="calories"
        class="elevation-2 mt-15"
        @update:items-per-page="handleCount"
    >
      <template v-slot:item.peer_id="{ item }">
        <div style="width: 210px; word-break: break-all;">
          <span>{{ item.peer_id }}</span>
        </div>
      </template>
      <template v-slot:item.CreatedAt="{ item }">
        <span>{{ new Date(item.CreatedAt * 1000).toLocaleString() }}</span>
      </template>
      <template v-slot:item.price="{ item }">
        <span>{{ item.price }} CSD</span>
      </template>
      <template v-slot:item.UpdatedAt="{ item }">
        <span>{{ new Date(item.UpdatedAt * 1000).toLocaleString() }}</span>
      </template>
      <template v-slot:top>
        <v-toolbar flat>
          <v-text-field
              v-model="query.name"
              append-icon="mdi-magnify"
              label="Search by name"
              class="pt-10"
              @keydown.enter="handleSearch"
          ></v-text-field>
          <v-spacer></v-spacer>
        </v-toolbar>
      </template>
      <template v-slot:no-data>No data</template>
    </v-data-table>
  </div>
</template>
<script>
import { fetchZeroAccessNodes } from '@/api'

export default {
  components: { },
  data: () => ({
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
      // { text: 'Public IP', align: 'start', value: 'ip' },
      { text: 'Host', value: 'addr' },
      // { text: 'Listen Port', value: 'port' },
      { text: 'Colo', value: 'colo' },
      { text: 'Price/Hour', value: 'price' },
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
    }
  }
}
</script>
