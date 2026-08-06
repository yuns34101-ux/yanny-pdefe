<template>
  <div class="page-container">
    <el-card>
      <div class="table-header">
        <span class="table-title">管理员列表</span>
        <el-button type="primary" v-if="hasPerm('admin:create')" @click="openDialog()">
          <el-icon><Plus /></el-icon> 新增管理员
        </el-button>
      </div>

      <el-table :data="list" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="real_name" label="姓名" width="120" />
        <el-table-column label="角色" min-width="200">
          <template #default="{ row }">
            <el-tag v-for="r in row.roles" :key="r.id" size="small" style="margin-right:4px">
              {{ r.name }}
            </el-tag>
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
        <el-table-column prop="created_at" label="创建时间" width="170" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const hasPerm = (code) => authStore.hasPermission(code)

const list = ref([])

onMounted(() => {
  ElMessage.info('管理员列表接口待后端实现')
})
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.table-title { font-size: 16px; font-weight: 600; }
</style>
