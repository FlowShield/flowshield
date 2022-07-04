import request from '@/utils/request'

export function fetchZeroAccessClients(params = {}) {
  return request({
    url: '/access/client',
    method: 'get',
    params
  })
}

export function postZeroAccessClient(data = {}) {
  return request({
    url: '/access/client',
    method: 'post',
    data
  })
}

export function putZeroAccessClient(data = {}) {
  return request({
    url: '/access/client',
    method: 'put',
    data
  })
}

export function deleteZeroAccessClient(id) {
  return request({
    url: `/access/client/${id}`,
    method: 'delete'
  })
}
