import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomeView')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login')
  },
  {
    path: '/clients',
    name: 'clients',
    component: () => import('@/views/clients')
  },
  {
    path: '/servers',
    name: 'servers',
    component: () => import('@/views/servers')
  },
  {
    path: '/resources',
    name: 'resources',
    component: () => import('@/views/resources')
  },
  {
    path: '/relay',
    name: 'relay',
    component: () => import('@/views/relay')
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
