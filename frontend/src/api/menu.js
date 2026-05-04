import request from '@/utils/request'

export function getMenuList() {
  return request({
    url: '/v1/menu/list',
    method: 'get'
  })
}

export function addMenu(data) {
  return request({
    url: '/v1/menu/add',
    method: 'post',
    data
  })
}

export function updateMenu(data) {
  return request({
    url: '/v1/menu/edit',
    method: 'post',
    data
  })
}

export function deleteMenu(data) {
  return request({
    url: '/v1/menu/delete',
    method: 'post',
    data
  })
}
