import { getUploadToken } from '@/api/video'
import { ElMessage } from 'element-plus'

// 上传目录控制：不同类型文件存入不同目录
const DIR_MAP = {
  entity_logo: 'images/entities',
  mp_icon: 'images/mp',
  video_cover: 'images/covers',
  video: 'videos',
}

/**
 * 上传文件到七牛云
 * @param {File} file - 文件对象
 * @param {string} dirType - 目录类型：entity_logo | mp_icon | video_cover | video
 * @returns {Promise<string>} - 返回 CDN URL
 */
export async function uploadToQiniu(file, dirType = 'video_cover') {
  if (!file) throw new Error('未选择文件')
  if (!DIR_MAP[dirType]) throw new Error('未知上传目录：' + dirType)

  // 1. 获取上传 Token
  const fileType = dirType === 'video' ? 'video' : 'image'
  const res = await getUploadToken(fileType)
  const { token, domain, upload_host } = res.data

  // 2. 生成路径：<dir>/YYYYMMDD/HHmmss_random.ext
  const ext = file.name.split('.').pop()
  const now = new Date()
  const date = now.toISOString().slice(0, 10).replace(/-/g, '')
  const ts = now.toISOString().replace(/[-:]/g, '').replace(/\..+/, '')
  const random = Math.random().toString(36).slice(2, 8)
  const key = `${DIR_MAP[dirType]}/${date}/${ts.slice(8)}_${random}.${ext}`

  // 3. 直传七牛
  const formData = new FormData()
  formData.append('token', token)
  formData.append('key', key)
  formData.append('file', file)

  const uploadRes = await fetch(upload_host, { method: 'POST', body: formData })
  if (!uploadRes.ok) throw new Error('上传失败：' + uploadRes.status)

  const result = await uploadRes.json()

  // 4. 构造 CDN URL（补全协议）
  const baseUrl = domain.startsWith('http') ? domain : 'https://' + domain
  return `${baseUrl}/${result.key}`
}

/**
 * 快捷上传图片（触发 file input）
 * @param {string} dirType - 目录类型
 * @returns {Promise<{url: string} | null>}
 */
export function pickAndUpload(dirType = 'video_cover') {
  return new Promise((resolve) => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = 'image/*'
    input.onchange = async (e) => {
      const file = e.target.files[0]
      if (!file) { resolve(null); return }
      try {
        const url = await uploadToQiniu(file, dirType)
        ElMessage.success('上传成功')
        resolve({ url })
      } catch (err) {
        ElMessage.error('上传失败：' + err.message)
        resolve(null)
      }
    }
    input.click()
  })
}
