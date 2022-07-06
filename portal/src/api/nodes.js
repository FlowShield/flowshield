import request from '@/utils/request'

// 获取node列表
export function fetchZeroAccessNodes(params = {}) {
  return request({
    url: '/node',
    method: 'get',
    params
  })
}
