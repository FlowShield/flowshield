<template>
  <div>
    <h3 class="font-weight-thin text-h3 mt-10">{{ total }} Relays</h3>
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
          <form-dialog @on-success="handleSearch"/>
        </v-toolbar>
      </template>
      <template v-slot:item.actions="{ item }">
        <v-icon small class="mr-2" @click="editItem(item)">mdi-pencil</v-icon>
        <v-icon small @click="deleteItem(item)">mdi-delete</v-icon>
      </template>
      <template v-slot:item.action="{ item }">
        <pem-dialog :data="item"/>
        <confirm-dialog @on-confirm="handleDelete" :data="item"/>
      </template>
      <template v-slot:no-data>No data</template>
    </v-data-table>
  </div>
</template>
<script>
import PemDialog from '@/components/pem-dialog'
import ConfirmDialog from '@/components/confirm-dialog'
import FormDialog from './components/form-dialog'
import { deleteZeroAccessRelay, fetchZeroAccessRelays } from '@/api'

export default {
  components: { PemDialog, FormDialog, ConfirmDialog },
  data: () => ({
    loading: false,
    query: {
      name: '',
      page: 1,
      limit_num: 15
    },
    tableHeaders: [
      { text: 'Name', align: 'start', value: 'name' },
      { text: 'Listen port', value: 'port' },
      { text: 'Expose port', value: 'out_port' },
      { text: 'Host', value: 'host' },
      { text: 'UUID', value: 'uuid' },
      { text: 'Created at', value: 'CreatedAt' },
      { text: 'Updated at', value: 'UpdatedAt' },
      { text: 'Action', value: 'action' }
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
      fetchZeroAccessRelays(this.query).then(res => {
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
    handleDelete(ref) {
      const item = ref.data

      deleteZeroAccessRelay(item.ID).then(_ => {
        ref.$close()
      }).finally(() => {
        this.handleSearch()
      })
    }
  }
}
</script>
