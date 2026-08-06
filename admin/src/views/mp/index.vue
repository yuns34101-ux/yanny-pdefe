<template>
  <div class="page-container">
    <el-card>
      <div class="table-header">
        <span class="table-title">小程序账号列表</span>
        <el-button type="primary" v-if="hasPerm('mp_account:create')" @click="openDialog()">
          <el-icon><Plus /></el-icon> 新增账号
        </el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="app_name" label="小程序名称" min-width="150">
          <template #default="{ row }">
            <div style="display:flex;align-items:center;gap:8px">
              <el-avatar v-if="row.app_icon" :size="32" :src="row.app_icon" />
              <span>{{ row.app_name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="app_id" label="AppID" width="220" />
        <el-table-column prop="description" label="备注" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170" />
      </el-table>
    </el-card>

    <!-- 新增弹窗 -->
    <el-dialog v-model="dialogVisible" title="新增小程序账号" width="550px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="小程序名称" prop="app_name">
          <el-input v-model="form.app_name" maxlength="100" />
        </el-form-item>
        <el-form-item label="AppID" prop="app_id">
          <el-input v-model="form.app_id" maxlength="64" placeholder="微信小程序 AppID" />
        </el-form-item>
        <el-form-item label="AppSecret" prop="app_secret">
          <el-input v-model="form.app_secret" type="password" show-password maxlength="128" placeholder="微信小程序 AppSecret" />
        </el-form-item>
        <el-form-item label="图标URL" prop="app_icon">
          <el-input v-model="form.app_icon" placeholder="http://..." />
        </el-form-item>
        <el-form-item label="备注" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="2" maxlength="300" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { listMpAccounts, createMpAccount } from '@/api/entity'
import { useAuthStore } from '@/store/auth'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const hasPerm = (code) => authStore.hasPermission(code)

const list = ref([])
const loading = ref(false)
const fetchList = async () => {
  loading.value = true
  try {
    const res = await listMpAccounts({ page: 1, page_size: 100 })
    list.value = res.data
  } finally { loading.value = false }
}

// 新增
const dialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const form = reactive({ app_id: '', app_secret: '', app_name: '', app_icon: '', description: '', status: 1 })
const rules = {
  app_id: [{ required: true, message: '请输入AppID' }],
  app_secret: [{ required: true, message: '请输入AppSecret' }],
  app_name: [{ required: true, message: '请输入小程序名称' }],
}

const openDialog = () => {
  Object.assign(form, { app_id: '', app_secret: '', app_name: '', app_icon: '', description: '', status: 1 })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    await createMpAccount(form)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    fetchList()
  } finally { submitting.value = false }
}

onMounted(fetchList)
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.table-title { font-size: 16px; font-weight: 600; }
</style>
