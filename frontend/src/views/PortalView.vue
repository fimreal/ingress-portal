<template>
  <div class="portal-view">
    <!-- 顶部栏 -->
    <header class="header">
      <div class="logo">
        <span class="logo-icon">🚀</span>
        <span class="logo-text">Ingress Portal</span>
      </div>

      <div class="search-box">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索服务..."
          @input="handleSearch"
        />
        <span class="search-icon">🔍</span>
      </div>

      <div class="actions">
        <button class="btn-refresh" @click="refresh" :disabled="loading">
          <span :class="{ 'spin': loading }">🔄</span>
        </button>

        <button
          class="btn-super"
          :class="{ 'active': isAuthenticated }"
          @click="handleSuperMode"
        >
          <span>{{ isAuthenticated ? '🔓' : '🔒' }}</span>
          <span>{{ isAuthenticated ? '管理' : 'Super' }}</span>
        </button>
      </div>
    </header>

    <!-- 分组筛选 -->
    <div class="group-filter" v-if="allGroups.length > 0">
      <button
        :class="{ active: selectedGroup === 'all' }"
        @click="selectedGroup = 'all'"
      >
        全部
      </button>
      <button
        v-for="group in allGroups"
        :key="group"
        :class="{ active: selectedGroup === group }"
        @click="selectedGroup = group"
      >
        {{ group }}
      </button>
    </div>

    <!-- 主内容区 -->
    <main class="content">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <p>正在加载...⏳</p>
      </div>

      <div v-else-if="filteredIngresses.length === 0" class="empty">
        <div class="empty-icon">🔍</div>
        <p>没有找到符合条件的 Ingress</p>
        <p class="empty-hint">{{ ingresses.length > 0 ? '尝试调整搜索条件' : '请检查 K8s 连接配置' }}</p>
      </div>

      <div v-else class="groups">
        <div
          v-for="group in filteredGroups"
          :key="group.name"
          class="group-section"
        >
          <h2 class="group-title">
            <span>{{ group.name }}</span>
            <span class="count">{{ group.ingresses.length }}</span>
          </h2>

          <div class="cards-grid">
            <IngressCard
              v-for="ingress in group.ingresses"
              :key="ingress.name + ingress.namespace"
              :ingress="ingress"
            />
          </div>
        </div>
      </div>
    </main>

    <!-- Super Mode 登录弹窗 -->
    <SuperModeLogin
      :show="showLoginModal"
      @close="showLoginModal = false"
      @success="onLoginSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useIngressStore } from '@/stores/ingress'
import { useAuthStore } from '@/stores/auth'
import IngressCard from '@/components/IngressCard.vue'
import SuperModeLogin from '@/components/SuperModeLogin.vue'

const router = useRouter()
const ingressStore = useIngressStore()
const authStore = useAuthStore()

const showLoginModal = ref(false)
const searchQuery = ref('')
const selectedGroup = ref('all')

// 计算属性
const loading = computed(() => ingressStore.loading)
const ingresses = computed(() => ingressStore.ingresses)
const isAuthenticated = computed(() => authStore.isAuthenticated)
const allGroups = computed(() => ingressStore.allGroups)

const filteredIngresses = computed(() => {
  let result = ingresses.value

  // 分组过滤
  if (selectedGroup.value !== 'all') {
    result = result.filter(ing => ing.group === selectedGroup.value)
  }

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(ing =>
      ing.name.toLowerCase().includes(query) ||
      ing.host.toLowerCase().includes(query) ||
      (ing.description?.toLowerCase().includes(query))
    )
  }

  return result
})

const filteredGroups = computed(() => {
  const groups = new Map<string, typeof ingresses.value>()

  for (const ing of filteredIngresses.value) {
    if (!ing.visible) continue

    const groupName = ing.group || ing.namespace || '未分组'
    if (!groups.has(groupName)) {
      groups.set(groupName, [])
    }
    groups.get(groupName)!.push(ing)
  }

  // 转换为数组并排序
  const result = []
  for (const [name, ingresses] of groups) {
    ingresses.sort((a, b) => (b.priority || 0) - (a.priority || 0) || a.name.localeCompare(b.name))
    result.push({ name, ingresses })
  }

  return result.sort((a, b) => a.name.localeCompare(b.name))
})

// 方法
function handleSearch() {
  // 搜索已在 computed 中自动处理
}

async function refresh() {
  await ingressStore.fetchIngresses()
}

function handleSuperMode() {
  if (isAuthenticated.value) {
    router.push('/admin')
  } else {
    showLoginModal.value = true
  }
}

function onLoginSuccess() {
  router.push('/admin')
}

onMounted(() => {
  ingressStore.fetchIngresses()
})
</script>

<style scoped>
.portal-view {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: white;
  border-bottom: 1px solid #e5e7eb;
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  gap: 2rem;
  position: sticky;
  top: 0;
  z-index: 100;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
}

.logo-icon {
  font-size: 1.5rem;
}

.search-box {
  flex: 1;
  max-width: 400px;
  position: relative;
}

.search-box input {
  width: 100%;
  padding: 0.5rem 2.5rem 0.5rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  transition: border-color 0.2s;
}

.search-box input:focus {
  outline: none;
  border-color: #3b82f6;
}

.search-icon {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  opacity: 0.5;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.actions button {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.btn-refresh {
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  color: #374151;
}

.btn-refresh:hover:not(:disabled) {
  background: #e5e7eb;
}

.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.btn-super {
  background: #3b82f6;
  border: 1px solid #3b82f6;
  color: white;
}

.btn-super:hover {
  background: #2563eb;
}

.btn-super.active {
  background: #10b981;
  border-color: #10b981;
}

.group-filter {
  display: flex;
  gap: 0.5rem;
  padding: 1rem 2rem;
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
  overflow-x: auto;
}

.group-filter button {
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  white-space: nowrap;
  cursor: pointer;
  transition: all 0.2s;
  background: white;
  border: 1px solid #d1d5db;
  color: #374151;
}

.group-filter button:hover {
  border-color: #9ca3af;
}

.group-filter button.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}

.loading, .empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem;
  text-align: center;
}

.spinner {
  width: 3rem;
  height: 3rem;
  border: 3px solid #e5e7eb;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty p {
  font-size: 1.125rem;
  color: #4b5563;
  margin-bottom: 0.5rem;
}

.empty-hint {
  font-size: 0.875rem;
  color: #9ca3af;
}

.groups {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.group-section {
  background: white;
  border-radius: 0.75rem;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.group-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
  font-size: 1.125rem;
  font-weight: 600;
  color: #1f2937;
}

.group-title .count {
  background: #e5e7eb;
  color: #6b7280;
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1rem;
}
