<template>
  <transition name="fade">
    <div v-if="show" class="modal-overlay" @click="handleBackdropClick">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>🔐 Super Mode 认证</h3>
          <button class="btn-close" @click="close">✕</button>
        </div>

        <div class="modal-body">
          <p class="hint">请输入 Token 或密码进入管理模式</p>

          <div class="input-group">
            <label>Token / 密码</label>
            <input
              v-model="credential"
              :type="showPassword ? 'text' : 'password'"
              placeholder="输入 Token 或密码..."
              @keyup.enter="submit"
              :disabled="loading"
            />
            <button
              type="button"
              class="btn-toggle"
              @click="showPassword = !showPassword"
            >
              {{ showPassword ? '🙈' : '👁' }}
            </button>
          </div>

          <div v-if="error" class="error">
            ⚠️ {{ error }}
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-secondary" @click="close" :disabled="loading">
            取消
          </button>
          <button
            class="btn-primary"
            @click="submit"
            :disabled="!credential || loading"
          >
            <span v-if="loading" class="loading-spinner">⏳</span>
            <span>{{ loading ? '验证中...' : '登录' }}</span>
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'

interface Props {
  show: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

const authStore = useAuthStore()

const credential = ref('')
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

watch(() => props.show, (newVal) => {
  if (!newVal) {
    // 重置状态
    credential.value = ''
    error.value = ''
    loading.value = false
  }
})

function close() {
  emit('close')
}

function handleBackdropClick() {
  if (!loading.value) {
    close()
  }
}

async function submit() {
  if (!credential.value) return

  loading.value = true
  error.value = ''

  try {
    // 尝试作为 Token 验证
    await authStore.login({ token: credential.value })
    emit('success')
    close()
  } catch (tokenErr) {
    // 尝试作为密码验证
    try {
      await authStore.login({ password: credential.value })
      emit('success')
      close()
    } catch (passErr) {
      error.value = 'Token 或密码错误'
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal {
  background: white;
  border-radius: 1rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1f2937;
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.25rem;
  color: #9ca3af;
  cursor: pointer;
  padding: 0.25rem;
  transition: color 0.2s;
}

.btn-close:hover {
  color: #4b5563;
}

.modal-body {
  padding: 1.25rem;
}

.hint {
  color: #6b7280;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.input-group {
  position: relative;
}

.input-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.375rem;
}

.input-group input {
  width: 100%;
  padding: 0.625rem 2.5rem 0.625rem 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.input-group input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.input-group input:disabled {
  background: #f3f4f6;
  cursor: not-allowed;
}

.btn-toggle {
  position: absolute;
  right: 0.5rem;
  bottom: 0.375rem;
  background: none;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  padding: 0.25rem;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.btn-toggle:hover {
  opacity: 1;
}

.error {
  margin-top: 0.75rem;
  padding: 0.625rem;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 0.375rem;
  color: #dc2626;
  font-size: 0.875rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1rem 1.25rem;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;
}

.modal-footer button {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary {
  background: white;
  border: 1px solid #d1d5db;
  color: #374151;
}

.btn-secondary:hover:not(:disabled) {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.btn-primary {
  background: #3b82f6;
  border: 1px solid #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
