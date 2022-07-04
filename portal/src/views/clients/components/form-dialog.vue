<template>
  <v-dialog
      v-model="dialog"
      persistent
      max-width="600px"
  >
    <template #activator="{ on, attrs }">
      <v-btn fab v-bind="attrs" color="primary" v-on="on">
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </template>
    <v-card>
      <v-card-title>
        <span class="text-h5">New client</span>
      </v-card-title>
      <v-card-text>
        <v-form v-model="valid">
          <v-container>
            <v-text-field
                v-model="form.name"
                label="Name"
                :rules="rule.name"
                required
            ></v-text-field>
            <v-text-field
                v-model.number="form.port"
                label="Listen port"
                :min="1"
                :max="65535"
                type="number"
                required
            ></v-text-field>
            <v-select
                v-model="form.server_id"
                label="Server"
                :items="serverItems"
                item-text="name"
                item-value="ID"
                :loading="loadingServer"
            >
            </v-select>
            <v-text-field
                v-model="form.target"
                :rules="rule.target"
                label="Resource"
                required
                hint="HOST:PORT"
            ></v-text-field>
            <v-text-field
                v-model.number="form.expire"
                label="Valid days"
                :min="1"
                type="number"
                required
            ></v-text-field>
          </v-container>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="purple darken-1" text @click="dialog = false">Close</v-btn>
        <v-btn
            color="purple darken-1"
            text
            @click="handleSubmit"
            :loading="submitting">Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import { fetchZeroAccessServers, postZeroAccessClient } from '@/api'

export default {
  data: () => ({
    dialog: false,
    valid: false,
    serverItems: [],
    loadingServer: false,
    submitting: false,
    form: {
      name: '',
      port: null,
      server_id: null,
      target: '',
      expire: 30
    },
    rule: {
      name: [
        v => !!v || 'Name is required'
      ],
      target: [
        v => !!v || 'Resource is required',
        v => v.includes(':') || 'IP:HOST'
      ]
    }
  }),
  created() {
    this.getServerOptions()
  },
  methods: {
    getServerOptions() {
      this.loadingServer = true
      // FIXME implement the lazy load server options
      fetchZeroAccessServers({ limit_num: 999 }).then(res => {
        this.serverItems = (res.data.list || [])
      }).finally(() => {
        this.loadingServer = false
      })
    },
    handleSubmit() {
      this.submitting = true

      const form = { ...this.form }
      const [host, port] = form.target.split(':')
      form.target = { host, port: +port }
      postZeroAccessClient(form).then(res => {
        this.$emit('on-success')
        this.$message.success()
        this.dialog = false
      }).finally(() => {
        this.submitting = false
      })
    }
  }
}
</script>
