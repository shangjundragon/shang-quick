import request from '@/utils/request'

export function getOperLogList(params) {
  return request({
    url: '/v1/operLog/list',
    method: 'get',
    params
  })
}
