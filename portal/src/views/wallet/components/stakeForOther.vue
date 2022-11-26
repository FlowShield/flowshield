<template>
  <v-dialog
      v-model="dialog"
      persistent
      max-width="600px"
  >
    <template #activator="{ on, attrs }">
      <v-btn x-large rounded v-bind="attrs" color="primary" v-on="on">
        <v-icon>mdi-wallet</v-icon>
        Stake for other accounts
      </v-btn>
    </template>
    <v-card>
      <v-card-title>
        <span class="text-h5">Stake</span>
      </v-card-title>
      <v-card-text>
        <v-form v-model="valid" ref="form">
          <v-container>
            <v-text-field
                v-model="address"
                label="Stake address"
                :rules="rule.address"
                required
            ></v-text-field>
          </v-container>
          <v-container>
            <v-radio-group
                v-model="stakeType"
                row
                :rules="rule.stakeType"
            >
              <v-radio
                  label="Fullnode"
                  value=1
              ></v-radio>
              <v-radio
                  label="Privoder"
                  value=2
              ></v-radio>
            </v-radio-group>
          </v-container>
          <v-container>
            <v-text-field
                v-model="amount"
                label="Stake amount"
                :rules="rule.amount"
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
import { stakeAmountForOther } from '@/utils/ethers'

export default {
  data: () => ({
    dialog: false,
    valid: false,
    serverItems: [],
    loadingServer: false,
    submitting: false,
    amount: 0,
    stakeType: 0,
    address: '',
    rule: {
      amount: [
        v => !!v || 'Amount is required. ' +
            ' Please enter stake amount !'
      ], stakeType: [
        v => !!v || 'type is required. ' +
            ' Please select stake type !'
      ], address: [
        v => !!v || 'type is required. ' +
            ' Please select stake type !'
      ]
    }
  }),
  methods: {
    async handleSubmit() {
      if (this.$refs.form.validate()) {
        this.submitting = true
        await stakeAmountForOther(this.stakeType, this.address, this.amount)
        this.submitting = false
      }
    }
  }
}
</script>
