import request from '@/utils/request'

export function fetchZeroAccessServers(params = {}) {
  return request({
    url: '/access/server',
    method: 'get',
    params
  })
}

export function postZeroAccessServer(data = {}) {
  return request({
    url: '/access/server',
    method: 'post',
    data
  })
}

export function putZeroAccessServer(data = {}) {
  return request({
    url: '/access/server',
    method: 'put',
    data
  })
}

export function deleteZeroAccessServer(id) {
  return request({
    url: `/access/server/${id}`,
    method: 'delete'
  })
}
