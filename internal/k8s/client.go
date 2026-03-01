package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/example/ingress-portal/pkg/models"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client K8s 客户端
type Client struct {
	clientset *kubernetes.Clientset
}

// NewClient 创建新的 K8s 客户端
func NewClient() (*Client, error) {
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("获取 K8s 配置失败: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 K8s 客户端失败: %w", err)
	}

	return &Client{clientset: clientset}, nil
}

// getConfig 获取 K8s 配置（集群内或本地）
func getConfig() (*rest.Config, error) {
	// 首先尝试集群内配置
	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}

	// 本地开发：使用 kubeconfig
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")

	// 允许通过环境变量指定
	if envPath := os.Getenv("KUBECONFIG"); envPath != "" {
		kubeconfig = envPath
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

// ListIngresses 列出所有 Ingress
func (c *Client) ListIngresses(ctx context.Context) ([]*models.IngressInfo, error) {
	ingressList, err := c.clientset.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("列出 Ingress 失败: %w", err)
	}

	var result []*models.IngressInfo
	for i := range ingressList.Items {
		ing := &ingressList.Items[i]
		info := c.parseIngress(ctx, ing)
		result = append(result, info)
	}

	return result, nil
}

// GetIngress 获取单个 Ingress
func (c *Client) GetIngress(ctx context.Context, namespace, name string) (*models.IngressInfo, error) {
	ing, err := c.clientset.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取 Ingress 失败: %w", err)
	}

	return c.parseIngress(ctx, ing), nil
}

// UpdateIngressVisibility 更新 Ingress 可见性
func (c *Client) UpdateIngressVisibility(ctx context.Context, namespace, name string, visible bool) error {
	ing, err := c.clientset.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("获取 Ingress 失败: %w", err)
	}

	// 更新 annotation
	if ing.Annotations == nil {
		ing.Annotations = make(map[string]string)
	}

	value := "false"
	if visible {
		value = "true"
	}
	ing.Annotations["portal.example.com/visible"] = value

	// 保存
	_, err = c.clientset.NetworkingV1().Ingresses(namespace).Update(ctx, ing, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新 Ingress 失败: %w", err)
	}

	return nil
}

// parseIngress 解析 Ingress 信息
func (c *Client) parseIngress(ctx context.Context, ing *networkingv1.Ingress) *models.IngressInfo {
	info := &models.IngressInfo{
		Name:         ing.Name,
		Namespace:    ing.Namespace,
		DiscoveredAt: ing.CreationTimestamp.Time,
		BackendStatus: models.HealthStatusUnknown,
	}

	// 解析 Annotations
	if ing.Annotations != nil {
		anno := ing.Annotations

		// 必需：visible
		if v, ok := anno["portal.example.com/visible"]; ok {
			info.Visible = strings.ToLower(v) == "true"
		} else {
			// 默认可见
			info.Visible = true
		}

		// 可选字段
		info.Group = anno["portal.example.com/group"]
		info.Description = anno["portal.example.com/description"]
		info.Team = anno["portal.example.com/team"]

		if p, err := strconv.Atoi(anno["portal.example.com/priority"]); err == nil {
			info.Priority = p
		}
	} else {
		// 没有 annotations，默认可见
		info.Visible = true
	}

	// 解析规则，获取 host 和 path
	if len(ing.Spec.Rules) > 0 {
		rule := ing.Spec.Rules[0]
		info.Host = rule.Host

		if len(rule.HTTP.Paths) > 0 {
			path := rule.HTTP.Paths[0]
			info.Path = path.Path

			// 解析后端 Service
			if path.Backend.Service != nil {
				info.Service = fmt.Sprintf("%s:%d",
					path.Backend.Service.Name,
					path.Backend.Service.Port.Number)
			}
		}
	}

	// 获取后端健康状态
	if info.Service != "" {
		info.BackendStatus = c.getBackendStatus(ctx, ing.Namespace, info.Service)
	}

	// Favicon URL
	if info.Host != "" {
		info.FaviconURL = fmt.Sprintf("https://%s/favicon.ico", info.Host)
	}

	return info
}

// getBackendStatus 获取后端健康状态
func (c *Client) getBackendStatus(ctx context.Context, namespace, serviceName string) models.HealthStatus {
	// 解析 service 名称和端口
	parts := strings.Split(serviceName, ":")
	if len(parts) != 2 {
		return models.HealthStatusUnknown
	}

	svcName := parts[0]

	// 获取 Endpoints
	endpoints, err := c.clientset.CoreV1().Endpoints(namespace).Get(ctx, svcName, metav1.GetOptions{})
	if err != nil {
		return models.HealthStatusUnknown
	}

	// 计算健康状态
	total := 0
	ready := 0

	for _, subset := range endpoints.Subsets {
		for _, addr := range subset.Addresses {
			total++
			if addr.TargetRef != nil {
				ready++
			}
		}
	}

	if total == 0 {
		return models.HealthStatusUnknown
	}

	if ready == total {
		return models.HealthStatusHealthy
	}

	if ready == 0 {
		return models.HealthStatusUnhealthy
	}

	return models.HealthStatusDegraded
}
