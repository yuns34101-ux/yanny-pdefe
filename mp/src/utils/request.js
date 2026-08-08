// API 基础配置
const BASE_URL = 'https://www-pdefe.yuangs.com/api/v1/mp'

// 请求封装
export function request(url, options = {}) {
  const token = uni.getStorageSync('mp_token')
  const header = { 'Content-Type': 'application/json', ...options.header }

  if (token) {
    header['Authorization'] = `Bearer ${token}`
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE_URL + url,
      method: options.method || 'GET',
      data: options.data,
      header,
      success(res) {
        const { data } = res
        if (data.code === 0) {
          resolve(data)
        } else if (data.code === 10002) {
          // 未登录，清除本地 token
          uni.removeStorageSync('mp_token')
          reject(new Error('请先登录'))
        } else {
          const err = new Error(data.message || '请求失败')
          err.code = data.code
          reject(err)
        }
      },
      fail(err) {
        reject(new Error('网络错误'))
      },
    })
  })
}

// GET 请求
export function get(url, params = {}) {
  const query = Object.keys(params)
    .filter(k => params[k] !== undefined && params[k] !== null && params[k] !== '')
    .map(k => `${k}=${encodeURIComponent(params[k])}`)
    .join('&')
  return request(query ? `${url}?${query}` : url, { method: 'GET' })
}

// POST 请求
export function post(url, data = {}) {
  return request(url, { method: 'POST', data })
}

// PUT 请求
export function put(url, data = {}) {
  return request(url, { method: 'PUT', data })
}

// 文件上传（后端代理中转，不直传七牛）
export function uploadFile(url, filePath, name = 'file') {
  const token = uni.getStorageSync('mp_token')
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: BASE_URL + url,
      filePath,
      name,
      header: { Authorization: 'Bearer ' + token },
      success(res) {
        try {
          const data = JSON.parse(res.data)
          if (data.code === 0) resolve(data.data)
          else reject(new Error(data.message || '上传失败'))
        } catch { reject(new Error('解析上传响应失败')) }
      },
      fail: reject,
    })
  })
}
