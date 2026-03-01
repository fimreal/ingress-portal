package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     []byte
	activeTokens  = make(map[string]*TokenInfo)
	tokensMu      sync.RWMutex
	globalConfig  Config
)

// Config 认证配置
type Config struct {
	Mode     string // "token" 或 "password"
	Password string // 当 Mode="password" 时使用
	TTL      int    // Token 有效期（小时）
}

// TokenInfo Token 信息
type TokenInfo struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	UsedCount int       `json:"usedCount"`
}

// TokenStatus Token 状态
type TokenStatus struct {
	Status        string     `json:"status"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	TotalRequests int        `json:"totalRequests"`
}

// JWTClaims JWT 声明
type JWTClaims struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

// Initialize 初始化认证系统
func Initialize(config Config) error {
	globalConfig = config

	// 生成随机 JWT 密钥
	jwtSecret = make([]byte, 32)
	if _, err := rand.Read(jwtSecret); err != nil {
		return fmt.Errorf("生成 JWT 密钥失败: %w", err)
	}

	return nil
}

// GenerateToken 生成新的 Super Mode Token
func GenerateToken(ttlHours int) (*TokenInfo, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}

	tokenStr := hex.EncodeToString(tokenBytes)
	now := time.Now()
	expiresAt := now.Add(time.Duration(ttlHours) * time.Hour)

	info := &TokenInfo{
		Token:     tokenStr,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}

	tokensMu.Lock()
	activeTokens[tokenStr] = info
	tokensMu.Unlock()

	return info, nil
}

// ValidateToken 验证 Super Mode Token
func ValidateToken(token string) (*TokenInfo, bool) {
	tokensMu.RLock()
	info, exists := activeTokens[token]
	tokensMu.RUnlock()

	if !exists {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(info.ExpiresAt) {
		// 删除过期 Token
		tokensMu.Lock()
		delete(activeTokens, token)
		tokensMu.Unlock()
		return nil, false
	}

	// 增加使用计数
	tokensMu.Lock()
	info.UsedCount++
	tokensMu.Unlock()

	return info, true
}

// GetActiveToken 获取当前激活的 Token（用于显示）
func GetActiveToken() (string, error) {
	tokensMu.RLock()
	defer tokensMu.RUnlock()

	// 返回最新创建的 Token
	var latest *TokenInfo
	for _, info := range activeTokens {
		if latest == nil || info.CreatedAt.After(latest.CreatedAt) {
			latest = info
		}
	}

	if latest == nil {
		return "", fmt.Errorf("没有激活的 Token")
	}

	return latest.Token, nil
}

// GetTokenStatus 获取 Token 状态
func GetTokenStatus() *TokenStatus {
	tokensMu.RLock()
	defer tokensMu.RUnlock()

	var totalRequests int
	var latestExpires *time.Time

	for _, info := range activeTokens {
		totalRequests += info.UsedCount
		if latestExpires == nil || info.ExpiresAt.After(*latestExpires) {
			t := info.ExpiresAt
			latestExpires = &t
		}
	}

	status := &TokenStatus{
		TotalRequests: totalRequests,
	}

	if len(activeTokens) == 0 {
		status.Status = "无激活 Token"
	} else if latestExpires != nil && time.Now().After(*latestExpires) {
		status.Status = "已过期"
	} else {
		status.Status = "正常"
		status.ExpiresAt = latestExpires
	}

	return status
}

// RevokeToken 吊销指定 Token
func RevokeToken(token string) {
	tokensMu.Lock()
	defer tokensMu.Unlock()

	delete(activeTokens, token)
}

// RevokeAllTokens 吊销所有 Token
func RevokeAllTokens() {
	tokensMu.Lock()
	defer tokensMu.Unlock()

	activeTokens = make(map[string]*TokenInfo)
}

// GenerateJWT 为用户生成 JWT Token
func GenerateJWT(userID string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		Type:   "super_mode",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(globalConfig.TTL) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT 验证 JWT Token
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的 token claims")
}

// VerifyPassword 验证密码
func VerifyPassword(password string) bool {
	if globalConfig.Mode != "password" {
		return false
	}
	return password == globalConfig.Password
}
