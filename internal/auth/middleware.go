package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SuperModeRequired Super Mode 认证中间件
func SuperModeRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证信息",
			})
			return
		}

		// Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "无效的认证格式",
			})
			return
		}

		// 验证 JWT
		claims, err := ValidateJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "无效的凭证: " + err.Error(),
			})
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// OptionalAuth 可选认证中间件（用于某些接口可以同时支持已登录和未登录用户）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		claims, err := ValidateJWT(parts[1])
		if err == nil {
			c.Set("userID", claims.UserID)
		}

		c.Next()
	}
}

// GetCurrentUser 从上下文获取当前用户 ID
func GetCurrentUser(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// IsSuperMode 检查当前是否处于 Super Mode（有用户登录）
func IsSuperMode(c *gin.Context) bool {
	_, exists := c.Get("userID")
	return exists
}