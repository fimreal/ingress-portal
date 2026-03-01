package models

import (
	"time"
)

// IngressInfo 表示一个 Ingress 的完整信息
type IngressInfo struct {
	// K8s 基本信息
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Host      string `json:"host"`
	Path      string `json:"path,omitempty"`
	Service   string `json:"service,omitempty"`

	// Annotations 解析
	Visible     bool   `json:"visible"`     // 必需
	Group       string `json:"group,omitempty"`
	Description string `json:"description,omitempty"`
	Team        string `json:"team,omitempty"`
	Priority    int    `json:"priority,omitempty"`

	// 运行时信息
	FaviconURL     string       `json:"faviconUrl,omitempty"`
	BackendStatus  HealthStatus `json:"backendStatus"`
	DiscoveredAt   time.Time    `json:"discoveredAt"`
	LastUpdatedAt  time.Time    `json:"lastUpdatedAt"`
}

// HealthStatus 后端健康状态
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusUnknown   HealthStatus = "unknown"
)

// IngressGroup 表示一个分组的 Ingress
type IngressGroup struct {
	Name      string         `json:"name"`
	Priority  int            `json:"priority"`
	Ingresses []*IngressInfo `json:"ingresses"`
}

// SuperModeToken Super Mode Token 信息
type SuperModeToken struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	UsedCount int       `json:"usedCount"`
}

// TokenStatus Token 状态响应
type TokenStatus struct {
	Status        string     `json:"status"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	TotalRequests int        `json:"totalRequests"`
}