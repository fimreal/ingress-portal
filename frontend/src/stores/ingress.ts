import { defineStore } from 'pinia'
import type { IngressInfo } from '@/types/ingress'

export const useIngressStore = defineStore('ingress', {
  state: () => ({
    ingresses: [] as IngressInfo[],
    loading: false,
    lastRefresh: null as number | null,
    searchQuery: '',
    selectedGroup: 'all'
  }),

  getters: {
    visibleIngresses: (state) => {
      return state.ingresses.filter(ing => ing.visible !== false)
    },

    filteredIngresses: (state) => {
      let result = state.ingresses

      // 分组过滤
      if (state.selectedGroup !== 'all') {
        result = result.filter(ing => ing.group === state.selectedGroup)
      }

      // 搜索过滤
      if (state.searchQuery) {
        const query = state.searchQuery.toLowerCase()
        result = result.filter(ing =>
          ing.name.toLowerCase().includes(query) ||
          ing.host.toLowerCase().includes(query) ||
          (ing.description && ing.description.toLowerCase().includes(query))
        )
      }

      return result
    },

    groupedIngresses: (state) => {
      const groups = new Map<string, IngressInfo[]>()

      for (const ing of state.ingresses) {
        // 使用 group 或 namespace 作为分组
        const groupName = ing.group || ing.namespace || '未分组'

        if (!groups.has(groupName)) {
          groups.set(groupName, [])
        }
        groups.get(groupName)!.push(ing)
      }

      // 转换为数组并按优先级排序
      const result: { name: string; ingresses: IngressInfo[] }[] = []
      for (const [name, ingresses] of groups) {
        // 组内按优先级和名称排序
        ingresses.sort((a, b) => {
          if (a.priority !== b.priority) {
            return (b.priority || 0) - (a.priority || 0)
          }
          return a.name.localeCompare(b.name)
        })

        result.push({ name, ingresses })
      }

      // 分组排序
      result.sort((a, b) => a.name.localeCompare(b.name))

      return result
    },

    allGroups: (state) => {
      const groups = new Set<string>()
      for (const ing of state.ingresses) {
        if (ing.group) {
          groups.add(ing.group)
        }
      }
      return Array.from(groups).sort()
    }
  },

  actions: {
    async fetchIngresses() {
      this.loading = true
      try {
        const response = await fetch('/api/ingresses')
        if (!response.ok) throw new Error('获取失败')
        const data = await response.json()
        this.ingresses = data.ingresses || []
        this.lastRefresh = Date.now()
      } catch (error) {
        console.error('获取 Ingress 失败:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async refreshIngresses() {
      try {
        const response = await fetch('/api/ingresses/refresh')
        if (!response.ok) throw new Error('刷新失败')
        await this.fetchIngresses()
      } catch (error) {
        console.error('刷新 Ingress 失败:', error)
        throw error
      }
    },

    setSearchQuery(query: string) {
      this.searchQuery = query
    },

    setSelectedGroup(group: string) {
      this.selectedGroup = group
    }
  }
})
