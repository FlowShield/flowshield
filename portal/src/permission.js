import router from './router'
import { EventBus } from '@/utils/event-bus'
import store from '@/store'

const whiteList = ['login', 'home', 'wallet', 'nodes', 'orders', 'resources'] // skip login

router.beforeEach(async(to, from, next) => {
  EventBus.$emit('app.loading', true)
  document.title = to?.meta?.title || 'FlowShield Portal'

  // goto login if needed
  try {
    await store.dispatch('getUserInfo')
  } catch (e) {
    if (!whiteList.includes(to.name)) {
      EventBus.$emit('app.message', 'Need login', 'warning')
      hideLoading()
      next({ name: 'login' })
    }
  }

  next()
})

router.afterEach(_ => {
  hideLoading()
})

function hideLoading() {
  window.setTimeout(() => {
    EventBus.$emit('app.loading', false)
  }, 100)
}
