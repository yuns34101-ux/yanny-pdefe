import request from '@/utils/request'

// 获取七牛上传 Token
export function getUploadToken(fileType) {
  return request.post('/upload/token', { file_type: fileType })
}

// 分类管理
export function listCategories(params) {
  return request.get('/categories', { params })
}

export function createCategory(data) {
  return request.post('/categories', data)
}

// 视频管理
export function listVideos(params) {
  return request.get('/videos', { params })
}

export function createVideo(data) {
  return request.post('/videos', data)
}

export function updateVideo(id, data) {
  return request.put(`/videos/${id}`, data)
}

export function deleteVideo(id) {
  return request.delete(`/videos/${id}`)
}
