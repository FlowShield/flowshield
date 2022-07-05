import request from '@/utils/request'

// 获取资源列表
export function fetchZeroAccessResources(params = {}) {
  return request({
    url: '/access/resource',
    method: 'get',
    params
  })
}

// 添加资源
export function postZeroAccessResource(data = {}) {
  return request({
    url: '/access/resource',
    method: 'post',
    data
  })
}

// 修改资源
export function putZeroAccessResource(data = {}) {
  return request({
    url: '/access/resource',
    method: 'put',
    data
  })
}

// 删除资源
export function deleteZeroAccessResource(id) {
  return request({
    url: `/access/resource/${id}`,
    method: 'delete'
  })
}
