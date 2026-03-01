package api

import (
	"net/http"
	"strings"

	"github.com/example/ingress-portal/internal/k8s"
	"github.com/example/ingress-portal/pkg/models"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, k8sClient *k8s.Client) {
	api := r.Group("/api")
	{
		api.GET("/ingresses", func(c *gin.Context) {
			handleListIngresses(c, k8sClient)
		})
		api.GET("/ingresses/refresh", func(c *gin.Context) {
			handleRefreshIngresses(c, k8sClient)
		})
		api.POST("/auth/super-mode", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"token":     "demo-token",
				"expiresAt": "2026-03-02T00:00:00Z",
			})
		})
	}

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
			return
		}
		c.File("web/dist/index.html")
	})
}

func handleListIngresses(c *gin.Context, client *k8s.Client) {
	ingresses, err := client.ListIngresses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 只返回可见的
	var filtered []*models.IngressInfo
	for _, ing := range ingresses {
		if ing.Visible {
			filtered = append(filtered, ing)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"ingresses": filtered,
		"count":     len(filtered),
	})
}

func handleRefreshIngresses(c *gin.Context, client *k8s.Client) {
	ingresses, err := client.ListIngresses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "refreshed",
		"count":   len(ingresses),
	})
}
