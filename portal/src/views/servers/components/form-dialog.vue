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
        <span class="text-h5">New server</span>
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
                v-model="form.host"
                label="Host"
                :rules="rule.host"
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
            <v-text-field
                v-model.number="form.out_port"
                label="Expose port"
                :min="1"
                :max="65535"
                type="number"
                required
            ></v-text-field>
            <v-select
                v-model="form.resource_id"
                label="Resources"
                :items="resourceItems"
                item-text="name"
                item-value="ID"
                :loading="loadingResource"
                multiple
            >
            </v-select>
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
import { fetchZeroAccessResources, postZeroAccessServer } from '@/api'

export default {
  data: () => ({
    dialog: false,
    valid: false,
    resourceItems: [],
    loadingResource: false,
    submitting: false,
    form: {
      host: '',
      name: '',
      port: null,
      out_port: null,
      resource_id: null
    },
    rule: {
      name: [
        v => !!v || 'Name is required'
      ],
      host: [
        v => !!v || 'Host is required'
      ],
      port: [
        v => !!v || 'Port is required'
      ],
      out_port: [
        v => !!v || 'Expose port is required'
      ],
      resource_id: [
        v => !!v || 'Resource is required'
      ]
    }
  }),
  created() {
    this.getResourceOptions()
  },
  methods: {
    getResourceOptions() {
      this.loadingResource = true
      // FIXME implement the lazy load resource options
      fetchZeroAccessResources({ limit_num: 999 }).then(res => {
        this.resourceItems = (res.data.list || [])
      }).finally(() => {
        this.loadingResource = false
      })
    },
    handleSubmit() {
      this.submitting = true

      const form = { ...this.form }
      form.resource_id = form.resource_id?.join(',')

      postZeroAccessServer(form).then(res => {
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
