<template>
  <div>
    <h3 class="font-weight-thin text-h3 mt-10">{{ total }} Resources</h3>
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
      <template v-slot:item.CreatedAt="{ item }">
        <span>{{ new Date(item.CreatedAt * 1000).toLocaleString() }}</span>
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
          <form-dialog v-if="isshowadd" @on-success="handleSearch"/>
        </v-toolbar>
      </template>
      <template v-slot:item.action="{ item }">
        <confirm-dialog v-if="isshowadd" @on-confirm="handleDelete" :data="item"/>
      </template>
      <template v-slot:no-data>No data</template>
    </v-data-table>
  </div>
</template>
<script>
import ConfirmDialog from '@/components/confirm-dialog'
import FormDialog from './components/form-dialog'
import { deleteZeroAccessResource, fetchZeroAccessResources } from '@/api'
import store from '@/store'

export default {
  components: { FormDialog, ConfirmDialog },
  data: () => ({
    isshowadd: false,
    loading: false,
    query: {
      name: '',
      page: 1,
      limit_num: 15
    },
    tableHeaders: [
      { text: 'Name', align: 'start', value: 'name' },
      { text: 'Type', value: 'type' },
      { text: 'Host', value: 'host' },
      { text: 'Port', value: 'port' },
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
      this.isshowadd = store.state.user.master
      this.loading = true
      fetchZeroAccessResources(this.query).then(res => {
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
      deleteZeroAccessResource(item.uuid).then(_ => {
        ref.$close()
      }).finally(() => {
        this.handleSearch()
      })
    }
  }
}
</script>
