import request from '@/utils/request'

// 上传前查重
export function checkMediaAsset(data) {
  return request.post('/upload/check', data)
}

// 上传成功后确认落库
export function confirmMediaAsset(data) {
  return request.post('/upload/confirm', data)
}
