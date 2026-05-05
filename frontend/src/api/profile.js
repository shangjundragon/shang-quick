import request from '@/utils/request'

export function updateProfile(data) {
  return request({
    url: '/v1/profile/update',
    method: 'post',
    data
  })
}
