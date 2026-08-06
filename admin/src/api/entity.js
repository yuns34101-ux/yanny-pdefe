import request from '@/utils/request'

// ========== 主体管理 CRUD ==========
export function listEntities(params) { return request.get('/entities', { params }) }
export function getEntity(id) { return request.get(`/entities/${id}`) }
export function createEntity(data) { return request.post('/entities', data) }
export function updateEntity(id, data) { return request.put(`/entities/${id}`, data) }
export function deleteEntity(id) { return request.delete(`/entities/${id}`) }

// ========== 小程序账号 CRUD ==========
export function listMpAccounts(params) { return request.get('/mp-accounts', { params }) }
export function createMpAccount(data) { return request.post('/mp-accounts', data) }
export function updateMpAccount(id, data) { return request.put(`/mp-accounts/${id}`, data) }

// ========== 绑定管理 ==========
export function bindEntityMp(data) { return request.post('/bindings', data) }
export function unbindEntityMp(data) { return request.delete('/bindings', { data }) }

// ========== CDN 配置 CRUD ==========
export function listCdnConfigs(params) { return request.get('/cdn', { params }) }
export function createCdnConfig(data) { return request.post('/cdn', data) }
export function updateCdnConfig(id, data) { return request.put(`/cdn/${id}`, data) }
export function deleteCdnConfig(id) { return request.delete(`/cdn/${id}`) }

// ========== 用户管理 ==========
export function listUsers(params) { return request.get('/users', { params }) }
export function updateUserStatus(id, status) { return request.put(`/users/${id}/status`, { status }) }

// ========== 管理员管理 CRUD ==========
export function listAdmins(params) { return request.get('/admins', { params }) }
export function createAdmin(data) { return request.post('/admins', data) }
export function updateAdmin(id, data) { return request.put(`/admins/${id}`, data) }
export function deleteAdmin(id) { return request.delete(`/admins/${id}`) }

// ========== 角色管理 ==========
export function listRoles() { return request.get('/roles') }
export function getRolePermissions(id) { return request.get(`/roles/${id}/permissions`) }
export function createRole(data) { return request.post('/roles', data) }
export function updateRole(id, data) { return request.put(`/roles/${id}`, data) }
