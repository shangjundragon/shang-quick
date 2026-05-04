import request from '@/utils/request'

export function getDeptList() {
  return request({
    url: '/v1/dept/list',
    method: 'get'
  })
}

export function addDept(data) {
  return request({
    url: '/v1/dept/add',
    method: 'post',
    data
  })
}

export function updateDept(data) {
  return request({
    url: '/v1/dept/edit',
    method: 'post',
    data
  })
}

export function deleteDept(data) {
  return request({
    url: '/v1/dept/delete',
    method: 'post',
    data
  })
}
