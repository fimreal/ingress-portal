export interface IngressInfo {
  name: string
  namespace: string
  host: string
  path?: string
  service?: string
  visible: boolean
  group?: string
  description?: string
  team?: string
  priority?: number
  faviconUrl?: string
  backendStatus: 'healthy' | 'degraded' | 'unhealthy' | 'unknown'
  discoveredAt: string
  lastUpdatedAt?: string
}

export interface IngressGroup {
  name: string
  priority: number
  ingresses: IngressInfo[]
}

export interface TokenInfo {
  token: string
  expiresAt: string
}
