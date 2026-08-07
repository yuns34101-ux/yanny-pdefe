<template>
  <div class="page-container">
    <el-card>
      <!-- 主体 + 搜索 -->
      <el-form :inline="true" :model="query">
        <el-form-item label="主体">
          <el-select v-model="query.entity_id" placeholder="选择主体" clearable style="width:200px" @change="onEntityChange">
            <el-option v-for="e in entities" :key="e.id" :label="e.name" :value="e.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="query.category_id" placeholder="全部分类" clearable style="width:160px">
            <el-option v-for="c in searchCategories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable style="width:120px">
            <el-option label="待审核" :value="0" />
            <el-option label="已发布" :value="1" />
            <el-option label="已下架" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="query.keyword" placeholder="标题" clearable style="width:180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">搜索</el-button>
        </el-form-item>
      </el-form>

      <div class="table-header">
        <span class="table-title">视频列表（共 {{ total }} 条）</span>
        <div>
          <el-button @click="openCategoryDialog">管理分类</el-button>
          <el-button type="primary" v-if="hasPerm('video:create')" @click="openDialog()">
            <el-icon><Plus /></el-icon> 新增视频
          </el-button>
        </div>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="封面" width="100">
          <template #default="{ row }">
            <el-image :src="row.cover_url" style="width:80px;height:60px;border-radius:4px" fit="cover" />
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="180" show-overflow-tooltip />
        <el-table-column label="时长" width="80">
          <template #default="{ row }">{{ formatDuration(row.duration) }}</template>
        </el-table-column>
        <el-table-column prop="view_count" label="播放" width="90" align="right" sortable />
        <el-table-column prop="like_count" label="点赞" width="80" align="right" sortable />
        <el-table-column prop="collect_count" label="收藏" width="80" align="right" />
        <el-table-column prop="comment_count" label="评论" width="80" align="right" />
        <el-table-column prop="status" label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="发布时间" width="170">
          <template #default="{ row }">{{ row.published_at || row.created_at }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openDialog(row)" v-if="hasPerm('video:edit')">编辑</el-button>
            <el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button link type="danger" size="small" v-if="hasPerm('video:delete')">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-footer">
        <el-pagination
          v-model:current-page="query.page" v-model:page-size="query.page_size"
          :total="total" :page-sizes="[10,20,50]"
          layout="total,sizes,prev,pager,next" @change="fetchList"
        />
      </div>
    </el-card>

    <!-- 新增/编辑视频 -->
    <el-dialog v-model="dialogVisible" :title="isEdit?'编辑视频':'新增视频'" width="650px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属主体" prop="entity_id">
              <el-select v-model="form.entity_id" placeholder="选择主体" style="width:100%" @change="onFormEntityChange">
                <el-option v-for="e in entities" :key="e.id" :label="e.name" :value="e.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属小程序" prop="mp_account_id">
              <el-input-number v-model="form.mp_account_id" :min="1" style="width:100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="分类" prop="category_id">
          <el-select v-model="form.category_id" placeholder="选择分类" style="width:100%">
            <el-option v-for="c in formCategories" :key="c.id" :label="c.name" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" maxlength="200" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" maxlength="1000" />
        </el-form-item>
        <el-form-item label="封面" prop="cover_url">
          <div style="display:flex;gap:8px">
            <el-input v-model="form.cover_url" placeholder="图片 URL 或点击上传" style="flex:1" />
            <el-button @click="triggerCoverUpload" :loading="uploadingCover">
              <el-icon><Upload /></el-icon> 上传
            </el-button>
          </div>
          <input ref="coverInput" type="file" accept="image/*" style="display:none" @change="handleCoverUpload" />
        </el-form-item>
        <el-form-item label="视频" prop="video_url">
          <div style="display:flex;gap:8px">
            <el-input v-model="form.video_url" placeholder="视频 URL 或点击上传" style="flex:1" />
            <el-button @click="triggerVideoUpload" :loading="uploadingVideo">
              <el-icon><Upload /></el-icon> 上传
            </el-button>
          </div>
          <input ref="videoInput" type="file" accept="video/*" style="display:none" @change="handleVideoUpload" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="时长(秒)">
              <el-input-number v-model="form.duration" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="状态">
              <el-select v-model="form.status">
                <el-option label="待审核" :value="0" />
                <el-option label="已发布" :value="1" />
                <el-option label="已下架" :value="2" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="推荐">
              <el-switch v-model="form.is_recommended" :active-value="1" :inactive-value="0" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="标签">
          <el-input v-model="form.tags" placeholder='JSON数组，如 ["搞笑","萌宠"]' />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible=false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>

    <!-- 分类管理弹窗 -->
    <el-dialog v-model="categoryVisible" title="分类管理" width="550px">
      <el-form :inline="true" style="margin-bottom:16px">
        <el-form-item label="主体">
          <el-select v-model="catEntityId" placeholder="选择主体" clearable style="width:220px" @change="loadDialogCategories">
            <el-option v-for="e in entities" :key="e.id" :label="e.name" :value="e.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <el-table :data="dialogCategories" size="small" v-if="catEntityId">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status===1?'success':'info'" size="small">{{ row.status===1?'启用':'禁用' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-else description="请先选择主体" :image-size="60" />
      <el-divider />
      <el-form :inline="true" :model="catForm" v-if="catEntityId">
        <el-form-item label="名称">
          <el-input v-model="catForm.name" placeholder="分类名称" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="catForm.sort_order" :min="0" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAddCategory">添加</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { listVideos, createVideo, updateVideo, deleteVideo, listCategories, createCategory, getUploadToken } from '@/api/video'
import { listEntities } from '@/api/entity'
import { useAuthStore } from '@/store/auth'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const hasPerm = (code) => authStore.hasPermission(code)

// 主体列表（全局）
const entities = ref([])

// 分类
const searchCategories = ref([])  // 搜索栏分类下拉
const dialogCategories = ref([])  // 分类管理弹窗表格
const formCategories = ref([])    // 视频表单中按 entity 过滤的分类
const categoryVisible = ref(false)
const catEntityId = ref(null)
const catForm = reactive({ name: '', sort_order: 0 })

const loadEntities = async () => {
  try { const res = await listEntities({ page_size: 200 }); entities.value = res.data || [] } catch { }
}

const loadSearchCategories = async (entityId) => {
  if (!entityId) { searchCategories.value = []; return }
  try { const res = await listCategories({ entity_id: entityId }); searchCategories.value = res.data || [] } catch { }
}

const loadDialogCategories = async () => {
  if (!catEntityId.value) { dialogCategories.value = []; return }
  try { const res = await listCategories({ entity_id: catEntityId.value }); dialogCategories.value = res.data || [] } catch { }
}

const loadFormCategories = async (entityId) => {
  if (!entityId) { formCategories.value = []; return }
  try { const res = await listCategories({ entity_id: entityId }); formCategories.value = res.data || [] } catch { }
}

const onEntityChange = () => {
  loadSearchCategories(query.entity_id)
  loadFormCategories(query.entity_id)
  query.category_id = null  // 切换主体时重置分类筛选
  fetchList()
}
const onFormEntityChange = (val) => { loadFormCategories(val) }

const openCategoryDialog = () => {
  if (!catEntityId.value && entities.value.length) {
    catEntityId.value = query.entity_id || entities.value[0]?.id
    loadDialogCategories()
  }
  categoryVisible.value = true
}

const handleAddCategory = async () => {
  if (!catForm.name) return ElMessage.warning('请输入分类名称')
  await createCategory({ entity_id: catEntityId.value, mp_account_id: 1, name: catForm.name, sort_order: catForm.sort_order })
  ElMessage.success('添加成功')
  catForm.name = ''
  loadDialogCategories()
  // 同时刷新搜索栏分类（如果是同一个主体）
  if (catEntityId.value === query.entity_id) {
    loadSearchCategories(query.entity_id)
    loadFormCategories(query.entity_id)
  }
}

// 视频列表
const list = ref([])
const total = ref(0)
const loading = ref(false)
const query = reactive({ entity_id: null, category_id: null, status: null, keyword: '', page: 1, page_size: 20 })

const fetchList = async () => {
  loading.value = true
  try {
    const params = { entity_id: query.entity_id, keyword: query.keyword, page: query.page, page_size: query.page_size }
    if (query.status || query.status === 0) params.status = query.status
    if (query.category_id) params.category_id = query.category_id
    const res = await listVideos(params)
    list.value = res.data; total.value = res.meta?.total || 0
  } finally { loading.value = false }
}

// 新增/编辑
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const editId = ref(null)
const form = reactive({ entity_id: null, mp_account_id: 1, category_id: null, title: '', description: '', cover_url: '', video_url: '', duration: 0, status: 0, is_recommended: 0, tags: '' })
const rules = {
  entity_id: [{ required: true, message: '选择主体', trigger: 'change' }],
  mp_account_id: [{ required: true, message: '必填' }],
  category_id: [{ required: true, message: '选择分类', trigger: 'change' }],
  cover_url: [{ required: true, message: '必填' }],
  video_url: [{ required: true, message: '必填' }],
}

// 文件上传
const coverInput = ref(null)
const videoInput = ref(null)
const uploadingCover = ref(false)
const uploadingVideo = ref(false)

const triggerCoverUpload = () => coverInput.value?.click()
const triggerVideoUpload = () => videoInput.value?.click()

const uploadToQiniu = async (file, fileType) => {
  const res = await getUploadToken(fileType)
  const { token, domain, upload_host } = res.data

  const ext = file.name.split('.').pop()
  const now = new Date()
  const ts = now.toISOString().replace(/[-:]/g, '').replace(/\..+/, '')
  const random = Math.random().toString(36).slice(2, 8)
  const key = `${fileType}s/${now.toISOString().slice(0, 10).replace(/-/g, '')}/${ts}_${random}.${ext}`

  const formData = new FormData()
  formData.append('token', token)
  formData.append('key', key)
  formData.append('file', file)

  const uploadRes = await fetch(upload_host, { method: 'POST', body: formData })
  if (!uploadRes.ok) throw new Error('上传失败')
  const result = await uploadRes.json()
  const baseUrl = domain.startsWith('http') ? domain : 'https://' + domain
  return `${baseUrl}/${result.key}`
}

const handleCoverUpload = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  uploadingCover.value = true
  try {
    form.cover_url = await uploadToQiniu(file, 'image')
    ElMessage.success('封面上传成功')
  } catch (err) {
    ElMessage.error('封面上传失败：' + err.message)
  } finally {
    uploadingCover.value = false
    if (coverInput.value) coverInput.value.value = ''
  }
}

const handleVideoUpload = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  uploadingVideo.value = true
  try {
    form.video_url = await uploadToQiniu(file, 'video')
    ElMessage.success('视频上传成功')
  } catch (err) {
    ElMessage.error('视频上传失败：' + err.message)
  } finally {
    uploadingVideo.value = false
    if (videoInput.value) videoInput.value.value = ''
  }
}

const openDialog = async (row) => {
  isEdit.value = !!row; editId.value = row?.id || null
  if (row) {
    Object.assign(form, { entity_id: row.entity_id, mp_account_id: row.mp_account_id, category_id: row.category_id, title: row.title, description: row.description, cover_url: row.cover_url, video_url: row.video_url, duration: row.duration, status: row.status, is_recommended: row.is_recommended, tags: row.tags })
    loadFormCategories(row.entity_id)
  } else {
    Object.assign(form, { entity_id: query.entity_id || null, mp_account_id: 1, category_id: null, title: '', description: '', cover_url: '', video_url: '', duration: 0, status: 0, is_recommended: 0, tags: '' })
    loadFormCategories(query.entity_id)
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  submitting.value = true
  try {
    if (isEdit.value) { await updateVideo(editId.value, form); ElMessage.success('更新成功') }
    else { await createVideo(form); ElMessage.success('创建成功') }
    dialogVisible.value = false; fetchList()
  } finally { submitting.value = false }
}

const handleDelete = async (id) => { await deleteVideo(id); ElMessage.success('删除成功'); fetchList() }

const statusType = (s) => ({ 0: 'warning', 1: 'success', 2: 'info' }[s] || 'info')
const statusText = (s) => ({ 0: '待审核', 1: '已发布', 2: '已下架' }[s] || '未知')
const formatDuration = (s) => { const m = Math.floor(s / 60), sec = s % 60; return `${m}:${String(sec).padStart(2, '0')}` }

onMounted(() => { loadEntities(); fetchList(); loadSearchCategories(query.entity_id) })
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin: 16px 0; }
.table-title { font-size: 16px; font-weight: 600; }
.table-footer { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
