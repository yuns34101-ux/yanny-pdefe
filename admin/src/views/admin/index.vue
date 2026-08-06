<template>
  <div class="page-container">
    <el-card>
      <div class="table-header">
        <span class="table-title">管理员列表</span>
        <el-button type="primary" @click="openDialog()">
          <el-icon><Plus /></el-icon> 新增管理员
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="username" label="用户名" width="140" />
        <el-table-column prop="real_name" label="姓名" width="120" />
        <el-table-column label="角色" min-width="200">
          <template #default="{ row }">
            <el-tag v-for="r in row.roles" :key="r.id" size="small" style="margin-right:4px">{{ r.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="管理主体" min-width="150">
          <template #default="{ row }">
            <template v-if="row.entity_ids">
              <el-tag v-for="eid in row.entity_ids" :key="eid" size="small" type="success" style="margin-right:4px">
                {{ entityMap[eid] || 'ID:' + eid }}
              </el-tag>
            </template>
            <el-tag v-else size="small" type="warning">全部</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status===1?'success':'danger'" size="small">{{ row.status===1?'正常':'禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openDialog(row)">编辑</el-button>
            <el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button link type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="visible" :title="isEdit?'编辑管理员':'新增管理员'" width="550px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="isEdit" maxlength="64" />
        </el-form-item>
        <el-form-item label="密码" :prop="isEdit?'':'password'">
          <el-input v-model="form.password" type="password" show-password :placeholder="isEdit?'留空不修改':'请输入密码'" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="form.real_name" maxlength="50" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="form.role_ids" multiple placeholder="选择角色" style="width:100%">
            <el-option v-for="r in roles" :key="r.id" :label="r.name" :value="r.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="管理主体">
          <el-select v-model="form.entity_ids" multiple placeholder="不选=全部" style="width:100%">
            <el-option v-for="e in entities" :key="e.id" :label="e.name" :value="e.id" />
          </el-select>
          <div style="font-size:12px;color:#909399;margin-top:4px">留空 = 管理所有主体（超级管理员）</div>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="visible=false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { listAdmins, createAdmin, updateAdmin, deleteAdmin, listRoles, listEntities } from '@/api/entity'
import { ElMessage } from 'element-plus'

const list = ref([])
const roles = ref([])
const entities = ref([])
const entityMap = ref({})
const loading = ref(false)
const visible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const editId = ref(null)
const form = reactive({ username: '', password: '', real_name: '', role_ids: [], entity_ids: [], status: 1 })
const rules = { username: [{ required: true, message: '请输入用户名' }], password: [{ required: true, message: '请输入密码' }] }

const fetchList = async () => {
  loading.value = true
  try { const res = await listAdmins(); list.value = res.data || [] } catch { } finally { loading.value = false }
}
const fetchMeta = async () => {
  try {
    const [rRes, eRes] = await Promise.all([listRoles(), listEntities({ page_size: 200 })])
    roles.value = rRes.data || []
    entities.value = eRes.data || []
    entityMap.value = {}; entities.value.forEach(e => { entityMap.value[e.id] = e.name })
  } catch { }
}
const openDialog = async (row) => {
  isEdit.value = !!row; editId.value = row?.id || null
  if (row) {
    Object.assign(form, { username: row.username, password: '', real_name: row.real_name, role_ids: (row.roles||[]).map(r=>r.id), entity_ids: row.entity_ids || [], status: row.status })
  } else {
    Object.assign(form, { username: '', password: '', real_name: '', role_ids: [], entity_ids: [], status: 1 })
  }
  visible.value = true
  nextTick(() => { if (formRef.value) formRef.value.clearValidate() })
}
const handleSubmit = async () => {
  if (isEdit.value && !form.password) rules.password = []
  else rules.password = [{ required: true, message: '请输入密码' }]
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) { await updateAdmin(editId.value, form); ElMessage.success('更新成功') }
    else { await createAdmin(form); ElMessage.success('创建成功') }
    visible.value = false; fetchList()
  } catch { } finally { submitting.value = false }
}
const handleDelete = async (id) => { await deleteAdmin(id); ElMessage.success('删除成功'); fetchList() }

onMounted(() => { fetchList(); fetchMeta() })
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.table-title { font-size: 16px; font-weight: 600; }
</style>
