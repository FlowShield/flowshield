import axios from 'axios'
// import qs from 'qs'
import { requestInterceptors } from './request-helper'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_URL,
  // withCredentials: false, // send cookies when cross-domain requests
  // transformRequest: [function(data, headers) {
  //   return qs.stringify(data)
  // }],
  // paramsSerializer: function(params) {
  //   return qs.stringify(params, { indices: false })
  // },
  timeout: 50000 // request timeout
})

requestInterceptors(service)

export default service
