import request from '@/utils/request'

export function getUserList(params) {
  return request({
    url: '/v1/user/list',
    method: 'get',
    params
  })
}

export function addUser(data) {
  return request({
    url: '/v1/user/add',
    method: 'post',
    data
  })
}

export function updateUser(data) {
  return request({
    url: '/v1/user/edit',
    method: 'post',
    data
  })
}

export function changeUserStatus(data) {
  return request({
    url: '/v1/user/changeStatus',
    method: 'post',
    data
  })
}

export function resetUserPwd(data) {
  return request({
    url: '/v1/user/resetPwd',
    method: 'post',
    data
  })
}

export function deleteUser(data) {
  return request({
    url: '/v1/user/delete',
    method: 'post',
    data
  })
}
