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
          reject(new Error(data.message || '请求失败'))
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
