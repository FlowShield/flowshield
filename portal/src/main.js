import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify'

// clip
import VueClipboard from 'vue-clipboard2'
VueClipboard.config.autoSetContainer = true // add this line
Vue.use(VueClipboard)

// Auth and permission
import './permission'

// global css
import './styles/index.scss'

Vue.config.productionTip = false

// mount function to Vue prototype, so you can use this.$message component
import { EventBus } from './utils/event-bus'

Vue.prototype.$message = {
  success: (msg = 'Success') => EventBus.$emit('app.message', msg, 'success'),
  error: (msg = 'Error') => EventBus.$emit('app.message', msg, 'error'),
  warning: (msg = 'Warning') => EventBus.$emit('app.message', msg, 'warning')
}

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount('#app')
