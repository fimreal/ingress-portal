import { defineStore } from 'pinia'
import type { TokenInfo } from '@/types/ingress'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('super_token') || '',
    expiresAt: localStorage.getItem('super_expires') || '',
    showLoginModal: false
  }),

  getters: {
    isAuthenticated: (state) => {
      if (!state.token) return false
      return new Date(state.expiresAt) > new Date()
    }
  },

  actions: {
    async login(credentials: { token?: string; password?: string }): Promise<void> {
      const response = await fetch('/api/auth/super-mode', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(credentials)
      })

      if (!response.ok) {
        const error = await response.json()
        throw new Error(error.error || '认证失败')
      }

      const data: TokenInfo = await response.json()
      this.token = data.token
      this.expiresAt = data.expiresAt

      localStorage.setItem('super_token', data.token)
      localStorage.setItem('super_expires', data.expiresAt)
      this.showLoginModal = false
    },

    logout() {
      fetch('/api/auth/super-mode', { method: 'DELETE' })
        .catch(() => {})

      this.token = ''
      this.expiresAt = ''
      localStorage.removeItem('super_token')
      localStorage.removeItem('super_expires')
    },

    openLoginModal() {
      this.showLoginModal = true
    },

    closeLoginModal() {
      this.showLoginModal = false
    }
  }
})
