<template>
  <div class="page-container">
    <el-card>
      <div class="table-header">
        <span class="table-title">角色与权限配置</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon> 新增角色
        </el-button>
      </div>

      <el-table :data="roles" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="code" label="编码" width="150" />
        <el-table-column prop="description" label="描述" min-width="220" show-overflow-tooltip />
        <el-table-column label="系统预置" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.is_system?'info':'success'" size="small">{{ row.is_system?'预置':'自定义' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="权限数" width="80" align="center">
          <template #default="{ row }">
            {{ (row.permissions || []).length }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewPerms(row)">查看权限</el-button>
            <el-button link type="primary" size="small" @click="openDialog(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑角色 -->
    <el-dialog v-model="visible" :title="isEdit?'编辑角色':'新增角色'" width="650px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" maxlength="50" />
        </el-form-item>
        <el-form-item label="编码" prop="code">
          <el-input v-model="form.code" :disabled="isEdit" maxlength="50" placeholder="英文标识，如 editor" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" maxlength="200" />
        </el-form-item>
        <el-form-item label="权限">
          <div style="width:100%">
            <el-checkbox v-model="checkAll" :indeterminate="isIndeterminate" @change="handleCheckAll">全选</el-checkbox>
            <div v-for="group in permGroups" :key="group.module" style="margin-top:12px">
              <el-divider content-position="left">
                <el-tag size="small">{{ group.label }}</el-tag>
              </el-divider>
              <el-checkbox-group v-model="form.permission_ids">
                <el-checkbox v-for="p in group.perms" :key="p.id" :value="p.id" :label="p.id">
                  {{ p.name }}
                </el-checkbox>
              </el-checkbox-group>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible=false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>

    <!-- 查看权限 -->
    <el-dialog v-model="permVisible" title="角色权限详情" width="600px">
      <el-table :data="permList" size="small">
        <el-table-column label="模块" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ moduleLabel(row.module) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="权限名称" min-width="180" />
        <el-table-column prop="code" label="编码" width="200" />
      </el-table>
      <template #footer><el-button @click="permVisible=false">关闭</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { listRoles, getRolePermissions, createRole, updateRole } from '@/api/entity'
import { ElMessage } from 'element-plus'

const roles = ref([])
const loading = ref(false)

// ===== CRUD =====
const visible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const editId = ref(null)
const form = reactive({ name: '', code: '', description: '', permission_ids: [] })
const rules = {
  name: [{ required: true, message: '请输入角色名称' }],
  code: [{ required: true, message: '请输入角色编码' }],
}

// 全选逻辑
const checkAll = computed({
  get: () => form.permission_ids.length === allPermIds.value.length,
  set: (v) => { form.permission_ids = v ? [...allPermIds.value] : [] },
})
const isIndeterminate = computed(() => form.permission_ids.length > 0 && form.permission_ids.length < allPermIds.value.length)

const allPermIds = ref([])
const permGroups = ref([])

const fetchList = async () => {
  loading.value = true
  try {
    const res = await listRoles()
    roles.value = res.data || []
  } catch { } finally { loading.value = false }
}

const openDialog = async (row) => {
  isEdit.value = !!row
  editId.value = row?.id || null

  // 加载权限列表
  try {
    const permRes = await listRoles() // 角色接口已返回 permissions
    // 从 super_admin 角色提取全量权限
    const superRole = (permRes.data || []).find(r => r.code === 'super_admin')
    if (superRole?.permissions) {
      allPermIds.value = superRole.permissions.map(p => p.id)
      // 按模块分组
      const groups = {}
      superRole.permissions.forEach(p => {
        if (!groups[p.module]) groups[p.module] = { module: p.module, label: moduleLabel(p.module), perms: [] }
        groups[p.module].perms.push(p)
      })
      permGroups.value = Object.values(groups)
    }
  } catch { }

  if (row) {
    Object.assign(form, {
      name: row.name, code: row.code, description: row.description,
      permission_ids: (row.permissions || []).map(p => p.id),
    })
  } else {
    Object.assign(form, { name: '', code: '', description: '', permission_ids: [] })
  }
  visible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      await updateRole(editId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createRole(form)
      ElMessage.success('创建成功')
    }
    visible.value = false
    fetchList()
  } catch { } finally { submitting.value = false }
}

const handleCheckAll = (val) => { form.permission_ids = val ? [...allPermIds.value] : [] }

// ===== 查看权限 =====
const permVisible = ref(false)
const permList = ref([])
const viewPerms = (row) => {
  permList.value = row.permissions || []
  permVisible.value = true
}

const moduleLabel = (m) => ({
  entity: '主体管理', mp_account: '小程序账号管理', cdn: 'CDN配置',
  video: '视频管理', user: '用户管理', analytics: '数据看板',
  admin: '管理员管理', role: '角色权限',
}[m] || m)

onMounted(fetchList)
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.table-title { font-size: 16px; font-weight: 600; }
</style>
