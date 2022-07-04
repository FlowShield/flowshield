<!--
Confirm dialog component

Props
  title:string  |  the title of the dialog
  btn-icon      |  the dialog placeholder button icon
  btn-text      |  the dialog placeholder button text
  data          |  attach some data to the component

Slots
  default    |  the dialog content

Events
  on-cancel  | the cancel button click event
  on-confirm | the confirm button click event

methods
  $show    | show the dialog
  $close   | close the dialog
-->
<template>
  <v-dialog v-model="dialog" scrollable max-width="600px" persistent>
    <template v-slot:activator="{ on, attrs }">
      <v-icon small v-bind="attrs" v-on="on">{{ btnIcon }}</v-icon>
      {{ btnText }}
    </template>
    <v-card>
      <v-card-title>{{ title }}</v-card-title>
      <v-card-text>
        <slot>Are you sure to delete?</slot>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text color="primary" @click="handleCancel">Cancel</v-btn>
        <v-btn color="primary" @click="handleConfirm" :loading="btnLoading">Confirm</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
export default {
  props: {
    title: {
      type: String,
      default: 'Confirm'
    },
    btnIcon: {
      type: String,
      default: 'mdi-delete'
    },
    btnText: {
      type: String,
      default: ''
    },
    data: {
      type: Object,
      default() {
        return null
      }
    }
  },
  data() {
    return {
      dialog: false,
      btnLoading: false
    }
  },
  methods: {
    handleCancel() {
      this.$emit('on-cancel', this)
      this.dialog = false
      this.btnLoading = false
    },
    handleConfirm() {
      this.$emit('on-confirm', this)
      this.btnLoading = true
    },
    $close() {
      this.btnLoading = false
      this.dialog = false
    },
    $show() {
      this.dialog = true
    }
  }
}
</script>
