import request from '@/utils/request'

export function getRoleList(params) {
  return request({
    url: '/v1/role/list',
    method: 'get',
    params
  })
}

export function addRole(data) {
  return request({
    url: '/v1/role/add',
    method: 'post',
    data
  })
}

export function updateRole(data) {
  return request({
    url: '/v1/role/edit',
    method: 'post',
    data
  })
}

export function deleteRole(data) {
  return request({
    url: '/v1/role/delete',
    method: 'post',
    data
  })
}

export function getRoleMenuIds(params) {
  return request({
    url: '/v1/role/menuIds',
    method: 'get',
    params
  })
}

export function assignRoleMenu(data) {
  return request({
    url: '/v1/role/assignMenu',
    method: 'post',
    data
  })
}
