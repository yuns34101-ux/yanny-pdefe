<template>
  <div class="page-container">
    <el-card>
      <div class="table-header">
        <span class="table-title">角色与权限配置</span>
        <el-button type="primary" v-if="hasPerm('role:create')">
          <el-icon><Plus /></el-icon> 新增角色
        </el-button>
      </div>

      <el-table :data="roles" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="code" label="编码" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column label="系统预置" width="100" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_system" type="info" size="small">预置</el-tag>
            <el-tag v-else type="success" size="small">自定义</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="权限数" width="100" align="center">
          <template #default="{ row }">
            {{ row.permissions?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewPerms(row)">查看权限</el-button>
            <el-button link type="primary" size="small" v-if="!row.is_system && hasPerm('role:edit')">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 权限查看弹窗 -->
    <el-dialog v-model="permVisible" title="角色权限详情" width="700px">
      <el-table :data="permList" size="small" v-if="permList.length">
        <el-table-column prop="module" label="模块" width="120">
          <template #default="{ row }">
            <el-tag>{{ moduleName(row.module) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="权限名称" min-width="180" />
        <el-table-column prop="code" label="编码" width="200" />
      </el-table>
      <template #footer>
        <el-button @click="permVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/store/auth'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const hasPerm = (code) => authStore.hasPermission(code)

// 预置角色数据（后续从 API 加载）
const roles = ref([
  { id: 1, name: '超级管理员', code: 'super_admin', description: '拥有全部权限', is_system: true },
  { id: 2, name: '管理员', code: 'admin', description: '日常运营管理，除系统配置外的全部权限', is_system: true },
  { id: 3, name: '内容编辑', code: 'editor', description: '仅视频内容管理与数据查看', is_system: true },
])

const permVisible = ref(false)
const permList = ref([])

const moduleName = (m) => ({
  entity: '主体管理', mp_account: '小程序账号', cdn: 'CDN配置',
  video: '视频管理', user: '用户管理', analytics: '数据分析',
  admin: '管理员', role: '角色权限',
}[m] || m)

const viewPerms = (row) => {
  permList.value = row.permissions || []
  permVisible.value = true
  if (!permList.value.length) {
    ElMessage.info('权限详情接口待后端实现')
  }
}
</script>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.table-title { font-size: 16px; font-weight: 600; }
</style>
