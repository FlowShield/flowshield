import request from '@/utils/request'

export function fetchUser() {
  return request({
    url: '/user/detail',
    method: 'get'
  })
}
