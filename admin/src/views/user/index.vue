<template>
  <div class="page-container">
    <el-card>
      <div class="table-header">
        <span class="table-title">用户管理</span>
        <span class="table-subtitle">按小程序维度查看用户</span>
      </div>

      <el-form :inline="true" :model="query">
        <el-form-item label="小程序">
          <el-select v-model="query.mp_account_id" placeholder="全部" clearable style="width: 200px">
            <el-option v-for="m in mpList" :key="m.id" :label="m.app_name" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="query.phone" placeholder="手机号搜索" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户信息" min-width="200">
          <template #default="{ row }">
            <div style="display:flex;align-items:center;gap:10px">
              <el-avatar :size="40" :src="row.avatar_url" />
              <div>
                <div>{{ row.nickname || '未设置昵称' }}</div>
                <div style="font-size:12px;color:#909399">{{ row.phone || '未绑定手机号' }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="gender" label="性别" width="80" align="center">
          <template #default="{ row }">
            {{ { 0: '未知', 1: '男', 2: '女' }[row.gender] || '未知' }}
          </template>
        </el-table-column>
        <el-table-column prop="province" label="地区" min-width="140">
          <template #default="{ row }">
            {{ row.province }} {{ row.city }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_login_at" label="最近登录" width="170" />
        <el-table-column prop="created_at" label="注册时间" width="170" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="toggleStatus(row)">
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
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
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { listMpAccounts } from '@/api/entity'
import { ElMessage } from 'element-plus'

// 小程序列表（供筛选）
const mpList = ref([])
const list = ref([])
const total = ref(0)
const query = reactive({ mp_account_id: null, phone: '', page: 1, page_size: 20 })

const fetchList = async () => {
  // TODO: 对接用户列表接口
  ElMessage.info('用户列表接口待后端实现')
}

const toggleStatus = (row) => {
  ElMessage.info('用户状态切换接口待后端实现')
}

onMounted(async () => {
  try { const res = await listMpAccounts({ page: 1, page_size: 100 }); mpList.value = res.data || [] } catch { }
})
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.table-header { margin-bottom: 16px; }
.table-title { font-size: 16px; font-weight: 600; margin-right: 12px; }
.table-subtitle { font-size: 13px; color: #909399; }
.table-footer { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
