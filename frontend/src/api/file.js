import request from '@/utils/request'

export function uploadFile(data) {
  return request({
    url: '/v1/file/upload',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function getFileList(params) {
  return request({
    url: '/v1/file/list',
    method: 'get',
    params
  })
}

export function deleteFile(data) {
  return request({
    url: '/v1/file/delete',
    method: 'post',
    data
  })
}

export function getFileConfig() {
  return request({
    url: '/v1/file/config',
    method: 'get'
  })
}
