import request from '@/utils/request'

export function login(username, password) {
  return request.post('/login', { username, password })
}

export function getAdminInfo() {
  return request.get('/info')
}

export function changePassword(oldPassword, newPassword) {
  return request.put('/password', { old_password: oldPassword, new_password: newPassword })
}
