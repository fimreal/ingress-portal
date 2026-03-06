<template>
  <div class="admin-view">
    <header class="header">
      <h1>管理面板</h1>
      <button @click="logout" class="btn-logout">退出登录</button>
    </header>

    <main class="content">
      <p>欢迎来到管理面板</p>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

function logout() {
  authStore.logout()
  router.push('/')
}
</script>

<style scoped>
.admin-view {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: white;
  border-bottom: 1px solid #e5e7eb;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h1 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #1f2937;
}

.btn-logout {
  padding: 0.5rem 1rem;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
}

.btn-logout:hover {
  background: #dc2626;
}

.content {
  flex: 1;
  padding: 2rem;
}
</style>
