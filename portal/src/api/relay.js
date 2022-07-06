import request from '@/utils/request'

// 获取relay列表
export function fetchZeroAccessRelays(params = {}) {
  return request({
    url: '/access/relay',
    method: 'get',
    params
  })
}

// 添加relay
export function postZeroAccessRelay(data = {}) {
  return request({
    url: '/access/relay',
    method: 'post',
    data
  })
}

// 修改relay
export function putZeroAccessRelay(data = {}) {
  return request({
    url: '/access/relay',
    method: 'put',
    data
  })
}

// 删除relay
export function deleteZeroAccessRelay(id) {
  return request({
    url: `/access/relay/${id}`,
    method: 'delete'
  })
}
