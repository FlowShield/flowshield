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
        <span class="text-h5">New order</span>
      </v-card-title>
      <v-card-text>
        <v-form v-model="valid" ref="form">
          <v-container>
            <v-text-field
                v-model="form.name"
                label="Name"
                :rules="rule.name"
                required
            ></v-text-field>
            <v-select
                v-model="form.peer_id"
                label="Server"
                :items="serverItems"
                item-text="showinfo"
                item-value="peer_id"
                :loading="loadingServer"
            >
            </v-select>
            <v-select
                v-model="form.resource_cid"
                label="Resource"
                :items="resourcesItems"
                item-text="name"
                item-value="cid"
                :loading="loadingResource"
            >
            </v-select>
            <v-text-field
                v-model.number="form.duration"
                label="Duration(hours)"
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
import { postZeroAccessClient } from '@/api'
import { fetchZeroAccessResources } from '@/api'
import { fetchZeroAccessNodes } from '@/api'

export default {
  data: () => ({
    dialog: false,
    valid: false,
    serverItems: [],
    resourcesItems: [],
    loadingServer: false,
    loadingResource: false,
    submitting: false,
    form: {
      name: '',
      peer_id: '',
      resource_cid: '',
      duration: 0
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
    this.getResourcesOptions()
  },
  methods: {
    getServerOptions() {
      this.loadingServer = true
      fetchZeroAccessNodes({ limit_num: 999 }).then(res => {
        this.serverItems = (res.data.list || [])
        for (let i = 0; i < this.serverItems.length; i++) {
          this.$set(this.serverItems[i], 'showinfo', this.serverItems[i].loc + ' --- ' + this.serverItems[i].addr + ' --- ' + this.serverItems[i].price + '   CSD/Hour')
        }
      }).finally(() => {
        this.loadingServer = false
      })
    },
    getResourcesOptions() {
      this.loadingResource = true
      fetchZeroAccessResources({ limit_num: 999 }).then(res => {
        this.resourcesItems = (res.data.list || [])
      }).finally(() => {
        this.loadingResource = false
      })
    },
    handleSubmit() {
      if (this.$refs.form.validate()) {
        this.submitting = true

        const form = { ...this.form }
        // const [host, port] = form.target.split(':')
        // form.target = { host, port: +port }
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
}
</script>
