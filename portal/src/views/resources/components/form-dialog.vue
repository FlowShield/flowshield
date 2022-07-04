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
        <span class="text-h5">New resource</span>
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
            <v-select
                v-model="form.type"
                label="Type"
                :items="['dns', 'cidr']"
            >
            </v-select>
            <v-text-field
                v-model="form.host"
                :rules="rule.host"
                label="Host"
                required
            ></v-text-field>
            <v-text-field
                v-model="form.port"
                label="Listen port"
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
import { postZeroAccessResource } from '@/api'

export default {
  data: () => ({
    dialog: false,
    valid: false,
    serverItems: [],
    loadingServer: false,
    submitting: false,
    form: {
      name: '',
      host: null,
      port: null,
      type: ''
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
      type: [
        v => !!v || 'Type is required'
      ]
    }
  }),
  methods: {
    handleSubmit() {
      this.submitting = true

      const form = { ...this.form }
      postZeroAccessResource(form).then(res => {
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
