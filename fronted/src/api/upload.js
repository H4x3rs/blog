import request from '../utils/request'

// 上传文件
export function uploadFile(file) {
  const formData = new FormData()
  formData.append('file', file)
  
  return request({
    url: '/upload/file',
    method: 'post',
    data: formData
  })
}

// 上传图片（用于编辑器）
export function uploadImage(file) {
  const formData = new FormData()
  formData.append('file', file)
  
  return request({
    url: '/upload/image',
    method: 'post',
    data: formData
  })
}

