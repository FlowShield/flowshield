<template>
  <v-dialog
      v-model="dialog"
      persistent
      max-width="600px"
  >
    <template #activator="{ on, attrs }">
      <v-btn x-large rounded v-bind="attrs" color="primary" v-on="on">
        <v-icon>mdi-wallet</v-icon>
        Change Your Wallet
      </v-btn>
    </template>
    <v-card>
      <v-card-title>
        <span class="text-h5">New wallet</span>
      </v-card-title>
      <v-card-text>
        <v-form v-model="valid" ref="form">
          <v-container>
            <v-text-field
                v-model="address"
                label="New wallet address"
                :rules="rule.address"
                required
            ></v-text-field>
          </v-container>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" text @click="dialog = false">Close</v-btn>
        <v-btn
            color="primary"
            text
            @click="handleSubmit"
            :loading="submitting">Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
import { changeWallet, getWallet } from '@/utils/ethers'
import store from '@/store'

export default {
  data: () => ({
    dialog: false,
    valid: false,
    serverItems: [],
    loadingServer: false,
    submitting: false,
    address: '',
    rule: {
      address: [
        v => !!v || 'Address is required. ' +
            ' Please enter a new wallet address !'
      ]
    }
  }),
  methods: {
    async handleSubmit() {
      if (this.$refs.form.validate()) {
        this.submitting = true
        await changeWallet(store.state.user.uuid, this.address)
        const address = await getWallet(store.state.user.uuid)
        if (address[0] === '0x0000000000000000000000000000000000000000') {
          this.$message.error('Change failed')
        } else {
          this.$emit('on-success')
          this.$message.success()
          this.dialog = false
        }
        this.submitting = false
      }
    }
  }
}
</script>
