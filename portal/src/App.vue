<template>
  <v-app>
    <v-progress-linear
        indeterminate
        fix
        :active="loading"
        height="1"
    ></v-progress-linear>
    <top-bar/>
    <v-main class="body-bg">
      <v-container>
        <router-view/>
      </v-container>
    </v-main>
    <v-snackbar
        app
        right
        top
        :color="level"
        v-model="snackbar"
    >
      {{ msg }}
      <template v-slot:action="{ attrs }">
        <v-btn
            text
            v-bind="attrs"
            @click="snackbar = false"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script>
import TopBar from '@/components/layout/top-bar'
import { EventBus } from '@/utils/event-bus'

export default {
  components: { TopBar },
  data: () => ({
    snackbar: false,
    loading: false,
    msg: '',
    level: ''
  }),
  mounted() {
    EventBus.$on('app.message', this.handleMessage)
    EventBus.$on('app.loading', this.handleAppLoading)
  },
  methods: {
    handleMessage(msg, level) {
      this.snackbar = true
      this.msg = msg
      this.level = level
    },
    handleAppLoading(flag) {
      this.loading = flag
    }
  }
}
</script>
<style lang="scss" scoped>
.body-bg {
  background: url("~@/assets/bg.png") no-repeat;
  background-size: cover;
}
</style>
