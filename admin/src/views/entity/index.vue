<template>
  <div class="page-container">
    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="query">
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="主体名称" clearable @keyup.enter="fetchList" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">搜索</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card class="table-card">
      <div class="table-header">
        <span class="table-title">主体列表</span>
        <el-button type="primary" @click="openDialog()" v-if="hasPerm('entity:create')">
          <el-icon><Plus /></el-icon> 新增主体
        </el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="主体名称" min-width="150">
          <template #default="{ row }">
            <div style="display:flex;align-items:center;gap:8px">
              <el-avatar v-if="row.logo_url" :size="32" :src="row.logo_url" />
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="contact_phone" label="联系电话" width="130" />
        <el-table-column prop="contact_email" label="邮箱" width="180" />
        <el-table-column prop="sort_order" label="排序" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openDialog(row)" v-if="hasPerm('entity:edit')">编辑</el-button>
            <el-button link type="primary" size="small" @click="manageBinding(row)">绑定小程序</el-button>
            <el-popconfirm title="确定删除该主体？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button link type="danger" size="small" v-if="hasPerm('entity:delete')">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-footer">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.page_size"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @change="fetchList"
        />
      </div>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑主体' : '新增主体'" width="600px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="主体名称" prop="name">
          <el-input v-model="form.name" maxlength="100" />
        </el-form-item>
        <el-form-item label="Logo URL" prop="logo_url">
          <el-input v-model="form.logo_url" placeholder="http://..." />
        </el-form-item>
        <el-form-item label="简介" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" maxlength="500" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="联系电话" prop="contact_phone">
              <el-input v-model="form.contact_phone" maxlength="20" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系邮箱" prop="contact_email">
              <el-input v-model="form.contact_email" maxlength="100" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="地址" prop="address">
          <el-input v-model="form.address" maxlength="300" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="排序" prop="sort_order">
              <el-input-number v-model="form.sort_order" :min="0" :max="9999" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-switch v-model="form.status" :active-value="1" :inactive-value="0" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>

    <!-- 绑定小程序弹窗 -->
    <el-dialog v-model="bindVisible" title="管理小程序绑定" width="550px">
      <div v-if="bindEntity" v-loading="bindLoading">
        <h4>{{ bindEntity.name }} - 已绑定小程序</h4>
        <el-table :data="bindings" style="margin-top:12px" v-if="bindings.length">
          <el-table-column label="小程序" min-width="180">
            <template #default="{ row }">
              {{ mpMap[row.mp_account_id] || '小程序 #' + row.mp_account_id }}
            </template>
          </el-table-column>
          <el-table-column prop="is_default" label="默认" width="80">
            <template #default="{ row }">
              <el-tag v-if="row.is_default" type="success" size="small">默认</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-popconfirm title="确定解绑？" @confirm="handleUnbind(row)">
                <template #reference>
                  <el-button link type="danger" size="small">解绑</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-else description="暂无绑定" :image-size="60" />

        <el-divider />
        <el-form :inline="true" :model="bindForm">
          <el-form-item label="小程序">
            <el-select v-model="bindForm.mp_account_id" placeholder="选择小程序" style="width:200px" filterable>
              <el-option v-for="m in mpList" :key="m.id" :label="m.app_name" :value="m.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="设为默认">
            <el-switch v-model="bindForm.is_default" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleBind">绑定</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { listEntities, createEntity, updateEntity, deleteEntity, bindEntityMp, unbindEntityMp, listMpAccounts, listEntityBindings } from '@/api/entity'
import { useAuthStore } from '@/store/auth'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const hasPerm = (code) => authStore.hasPermission(code)

// 列表
const list = ref([])
const total = ref(0)
const loading = ref(false)
const query = reactive({ keyword: '', status: null, page: 1, page_size: 20 })

const fetchList = async () => {
  loading.value = true
  try {
    const res = await listEntities(query)
    list.value = res.data
    total.value = res.meta?.total || 0
  } finally { loading.value = false }
}
const resetQuery = () => {
  query.keyword = ''
  query.status = null
  query.page = 1
  fetchList()
}

// 新增/编辑
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const form = reactive({ name: '', logo_url: '', description: '', contact_phone: '', contact_email: '', address: '', sort_order: 0, status: 1 })
const editId = ref(null)

const rules = {
  name: [{ required: true, message: '请输入主体名称', trigger: 'blur' }],
}

const openDialog = (row) => {
  isEdit.value = !!row
  editId.value = row?.id || null
  if (row) {
    Object.assign(form, row)
  } else {
    Object.assign(form, { name: '', logo_url: '', description: '', contact_phone: '', contact_email: '', address: '', sort_order: 0, status: 1 })
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) {
      await updateEntity(editId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createEntity(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchList()
  } finally { submitting.value = false }
}

const handleDelete = async (id) => {
  await deleteEntity(id)
  ElMessage.success('删除成功')
  fetchList()
}

// 绑定管理
const bindVisible = ref(false)
const bindLoading = ref(false)
const bindEntity = ref(null)
const bindings = ref([])
const mpList = ref([])
const mpMap = ref({})
const bindForm = reactive({ mp_account_id: null, is_default: false })

const manageBinding = async (row) => {
  bindEntity.value = row
  bindForm.mp_account_id = null
  bindForm.is_default = false
  bindLoading.value = true
  try {
    // 加载小程序账号列表
    const mpRes = await listMpAccounts({ page_size: 200 })
    mpList.value = mpRes.data || []
    mpMap.value = {}
    mpList.value.forEach(m => { mpMap.value[m.id] = m.app_name })

    // 加载已绑定列表
    const bindRes = await listEntityBindings(row.id)
    bindings.value = bindRes.data || []
  } catch { }
  bindLoading.value = false
  bindVisible.value = true
}

const handleBind = async () => {
  if (!bindForm.mp_account_id) return ElMessage.warning('请选择小程序')
  await bindEntityMp({ entity_id: bindEntity.value.id, mp_account_id: bindForm.mp_account_id, is_default: bindForm.is_default ? 1 : 0 })
  ElMessage.success('绑定成功')
  // 刷新绑定列表
  bindings.value.push({ mp_account_id: bindForm.mp_account_id, is_default: bindForm.is_default })
  bindForm.mp_account_id = null
  bindForm.is_default = false
}

const handleUnbind = async (row) => {
  await unbindEntityMp({ entity_id: bindEntity.value.id, mp_account_id: row.mp_account_id })
  ElMessage.success('解绑成功')
  bindings.value = bindings.value.filter(b => b.mp_account_id !== row.mp_account_id)
}

onMounted(fetchList)
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.search-card :deep(.el-card__body) { padding-bottom: 0; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.table-title { font-size: 16px; font-weight: 600; }
.table-footer { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
