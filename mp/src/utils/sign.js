// 埋点签名工具 — 防止伪造上报数据
// 算法：sign = MD5(video_id:timestamp:nonce:jwt_secret)

const SECRET = 'yanny-jwt-secret-change-in-production'

function generateNonce() {
  return Math.random().toString(36).substring(2, 10) + Date.now().toString(36)
}

// ========== MD5 (RFC 1321, cmn 风格) ==========

function md5(s) {
  const bytes = []
  for (let i = 0; i < s.length; i++) {
    bytes.push(s.charCodeAt(i) & 0xFF)
  }

  const msgLen = bytes.length
  // 追加 padding: 0x80 + zeros
  bytes.push(0x80)
  while ((bytes.length % 64) !== 56) {
    bytes.push(0x00)
  }
  // 追加 64-bit message length in bits (little-endian)
  const bits = msgLen * 8
  for (let i = 0; i < 8; i++) {
    bytes.push((bits >>> (i * 8)) & 0xFF)
  }

  // 将 bytes 转为 32-bit words (little-endian)
  const words = new Array(bytes.length >> 2)
  for (let i = 0; i < bytes.length; i += 4) {
    words[i >> 2] = bytes[i] | (bytes[i + 1] << 8) | (bytes[i + 2] << 16) | (bytes[i + 3] << 24)
  }

  // MD5 常量
  let a = 0x67452301, b = 0xEFCDAB89, c = 0x98BADCFE, d = 0x10325476

  function add(x, y) { return (x + y) & 0xFFFFFFFF }
  function rotl(x, n) { return ((x << n) | (x >>> (32 - n))) & 0xFFFFFFFF }
  function F(x, y, z) { return (x & y) | (~x & z) }
  function G(x, y, z) { return (x & z) | (y & ~z) }
  function H(x, y, z) { return x ^ y ^ z }
  function I(x, y, z) { return y ^ (x | ~z) }

  const S = [
    7,12,17,22, 7,12,17,22, 7,12,17,22, 7,12,17,22,
    5, 9,14,20, 5, 9,14,20, 5, 9,14,20, 5, 9,14,20,
    4,11,16,23, 4,11,16,23, 4,11,16,23, 4,11,16,23,
    6,10,15,21, 6,10,15,21, 6,10,15,21, 6,10,15,21,
  ]
  const K = [
    0xD76AA478,0xE8C7B756,0x242070DB,0xC1BDCEEE,0xF57C0FAF,0x4787C62A,0xA8304613,0xFD469501,
    0x698098D8,0x8B44F7AF,0xFFFF5BB1,0x895CD7BE,0x6B901122,0xFD987193,0xA679438E,0x49B40821,
    0xF61E2562,0xC040B340,0x265E5A51,0xE9B6C7AA,0xD62F105D,0x02441453,0xD8A1E681,0xE7D3FBC8,
    0x21E1CDE6,0xC33707D6,0xF4D50D87,0x455A14ED,0xA9E3E905,0xFCEFA3F8,0x676F02D9,0x8D2A4C8A,
    0xFFFA3942,0x8771F681,0x6D9D6122,0xFDE5380C,0xA4BEEA44,0x4BDECFA9,0xF6BB4B60,0xBEBFBC70,
    0x289B7EC6,0xEAA127FA,0xD4EF3085,0x04881D05,0xD9D4D039,0xE6DB99E5,0x1FA27CF8,0xC4AC5665,
    0xF4292244,0x432AFF97,0xAB9423A7,0xFC93A039,0x655B59C3,0x8F0CCC92,0xFFEFF47D,0x85845DD1,
    0x6FA87E4F,0xFE2CE6E0,0xA3014314,0x4E0811A1,0xF7537E82,0xBD3AF235,0x2AD7D2BB,0xEB86D391,
  ]

  for (let bi = 0; bi < words.length; bi += 16) {
    let aa = a, bb = b, cc = c, dd = d

    for (let i = 0; i < 64; i++) {
      let f, g
      if (i < 16)       { f = F(b, c, d); g = i }
      else if (i < 32)  { f = G(b, c, d); g = (5 * i + 1) % 16 }
      else if (i < 48)  { f = H(b, c, d); g = (3 * i + 5) % 16 }
      else              { f = I(b, c, d); g = (7 * i) % 16 }
      f = add(f, add(a, add(K[i], words[bi + g])))
      a = d; d = c; c = b; b = add(b, rotl(f, S[i]))
    }

    a = add(a, aa); b = add(b, bb); c = add(c, cc); d = add(d, dd)
  }

  // 输出 32 位 hex（little-endian bytes）
  function toHex(v) {
    let h = ''
    for (let i = 0; i < 4; i++) {
      const b = (v >>> (i * 8)) & 0xFF
      h += ('0' + b.toString(16)).slice(-2)
    }
    return h
  }
  return toHex(a) + toHex(b) + toHex(c) + toHex(d)
}

// ========== 签名 ==========

export function generateSign(id, timestamp) {
  const nonce = generateNonce()
  const raw = `${id}:${timestamp}:${nonce}:${SECRET}`
  const sign = md5(raw)
  return { timestamp, nonce, sign }
}

export function createTrackPayload(baseData) {
  const timestamp = Date.now()
  const id = baseData.video_id || baseData.target_id || 0
  const { sign, nonce } = generateSign(id, timestamp)
  return { ...baseData, timestamp, nonce, sign }
}
