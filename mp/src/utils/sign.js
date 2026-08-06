// 埋点签名工具 — 防止伪造上报数据
// 算法：sign = MD5(id + timestamp + nonce + secret)
// secret 与服务端 config.yaml jwt.secret 保持一致

const SECRET = 'yanny-jwt-secret-change-in-production'

// 生成随机 nonce
function generateNonce() {
  return Math.random().toString(36).substring(2, 10) + Date.now().toString(36)
}

// MD5 简化实现（小程序环境可用）
function md5(str) {
  // 使用 uni-app 内置或微信 API
  // 简化版：实际项目用 crypto-js 或微信原生 md5
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash |= 0
  }
  return Math.abs(hash).toString(16).padStart(8, '0')
}

// 生成埋点签名
export function generateSign(id, timestamp) {
  const nonce = generateNonce()
  const raw = `${id}:${timestamp}:${nonce}:${SECRET}`
  const sign = md5(raw)
  return { timestamp, nonce, sign }
}

// 创建带签名的埋点数据
export function createTrackPayload(baseData) {
  const timestamp = Date.now()
  const { sign, nonce } = generateSign(baseData.video_id || baseData.target_id || 0, timestamp)
  return {
    ...baseData,
    timestamp,
    nonce,
    sign,
  }
}
