import Vue from 'vue'
import Vuex from 'vuex'
import { fetchUser } from '@/api'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    user: {
      avatar: '',
      email: ''
    }
  }, getters: {}, mutations: {
    SET_USER(state, info) {
      state.user.avatar = info.avatar_url
      state.user.email = info.email
    }
  }, actions: {
    getUserInfo({ commit, state }) {
      return new Promise((resolve, reject) => {
        if (state.user.email) {
          resolve()
          return
        }

        fetchUser().then(res => {
          commit('SET_USER', res.data)
          resolve()
        }).catch(error => {
          reject(error)
        })
      })
    }
  }, modules: {}
})
