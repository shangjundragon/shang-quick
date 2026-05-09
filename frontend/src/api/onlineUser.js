import request from '@/utils/request'

export function getOnlineUserList() {
  return request({
    url: '/v1/onlineUser/list',
    method: 'get'
  })
}

export function kickOnlineUser(data) {
  return request({
    url: '/v1/onlineUser/kick',
    method: 'post',
    data
  })
}
