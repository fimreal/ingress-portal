<template>
  <a
    :href="ingressUrl"
    target="_blank"
    class="ingress-card"
    :class="cardClass"
  >
    <div class="card-header">
      <img
        v-if="ingress.faviconUrl"
        :src="ingress.faviconUrl"
        class="favicon"
        @error="onFaviconError"
        alt=""
      />
      <div v-else class="favicon-placeholder">🌐</div>

      <div
        class="status-badge"
        :class="'status-' + ingress.backendStatus"
        :title="statusText"
      />
    </div>

    <h3 class="ingress-name">{{ ingress.name }}</h3>

    <p v-if="ingress.description" class="description">
      {{ ingress.description }}
    </p>
    <p v-else class="host-text">{{ ingress.host }}</p>

    <div class="card-footer">
      <span v-if="ingress.team" class="team">{{ ingress.team }}</span>
      <span class="namespace">{{ ingress.namespace }}</span>
    </div>
  </a>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { IngressInfo } from '@/types/ingress'

interface Props {
  ingress: IngressInfo
}

const props = defineProps<Props>()

const faviconFailed = ref(false)

const ingressUrl = computed(() => {
  const protocol = 'https://'
  const path = props.ingress.path || ''
  return `${protocol}${props.ingress.host}${path}`
})

const cardClass = computed(() => ({
  'status-healthy': props.ingress.backendStatus === 'healthy',
  'status-degraded': props.ingress.backendStatus === 'degraded',
  'status-unhealthy': props.ingress.backendStatus === 'unhealthy',
}))

const statusText = computed(() => {
  const statusMap: Record<string, string> = {
    healthy: '后端健康',
    degraded: '后端部分可用',
    unhealthy: '后端异常',
    unknown: '状态未知'
  }
  return statusMap[props.ingress.backendStatus] || '状态未知'
})

function onFaviconError() {
  faviconFailed.value = true
}
</script>

<style scoped>
.ingress-card {
  display: block;
  background: white;
  border-radius: 0.75rem;
  padding: 1.25rem;
  text-decoration: none;
  color: inherit;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
  border: 2px solid transparent;
}

.ingress-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.ingress-card.status-healthy:hover {
  border-color: #10b981;
}

.ingress-card.status-degraded:hover {
  border-color: #f59e0b;
}

.ingress-card.status-unhealthy:hover {
  border-color: #ef4444;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

.favicon, .favicon-placeholder {
  width: 2rem;
  height: 2rem;
  border-radius: 0.375rem;
  object-fit: contain;
}

.favicon-placeholder {
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
}

.status-badge {
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 50%;
  margin-left: auto;
}

.status-badge.status-healthy {
  background: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.3);
}

.status-badge.status-degraded {
  background: #f59e0b;
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.3);
}

.status-badge.status-unhealthy {
  background: #ef4444;
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.3);
}

.status-badge.status-unknown {
  background: #9ca3af;
}

.ingress-name {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 0.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.description {
  font-size: 0.875rem;
  color: #6b7280;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 0.75rem;
  min-height: 2.5rem;
}

.host-text {
  font-size: 0.75rem;
  color: #9ca3af;
  margin-bottom: 0.75rem;
  min-height: 2.5rem;
}

.card-footer {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.team, .namespace {
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  border-radius: 0.25rem;
  font-weight: 500;
}

.team {
  background: #dbeafe;
  color: #1e40af;
}

.namespace {
  background: #f3f4f6;
  color: #4b5563;
}
</style>
