package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/ingress-portal/internal/api"
	"github.com/example/ingress-portal/internal/auth"
	"github.com/example/ingress-portal/internal/k8s"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	port        string
	superMode   string // token 或 password
	superPass   string
	tokenTTL    int    // 小时

	rootCmd = &cobra.Command{
		Use:   "ingress-portal",
		Short: "K8s Ingress 自动发现与导航门户",
		Long: `自动发现 K8s 集群中的 Ingress 资源，
提供统一的 Web 入口导航服务，支持 Super Mode 管理可见性。`,
		Run: runServer,
	}

	tokenCmd = &cobra.Command{
		Use:   "token",
		Short: "管理 Super Mode Token",
	}

	tokenGenCmd = &cobra.Command{
		Use:   "generate",
		Short: "生成新的 Super Mode Token",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := auth.GenerateToken(tokenTTL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "生成失败: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Token: %s\n", token.Token)
			fmt.Printf("有效期: %d 小时 (至 %s)\n", tokenTTL, token.ExpiresAt.Format("2006-01-02 15:04:05"))
		},
	}

	tokenStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "查看当前 Token 状态",
		Run: func(cmd *cobra.Command, args []string) {
			status := auth.GetTokenStatus()
			fmt.Printf("Token 状态: %s\n", status.Status)
			if status.ExpiresAt != nil {
				fmt.Printf("过期时间: %s\n", status.ExpiresAt.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("总请求数: %d\n", status.TotalRequests)
		},
	}

	tokenRevokeCmd = &cobra.Command{
		Use:   "revoke [token]",
		Short: "吊销指定 Token 或所有 Token",
		Run: func(cmd *cobra.Command, args []string) {
			all, _ := cmd.Flags().GetBool("all")
			if all {
				auth.RevokeAllTokens()
				fmt.Println("所有 Token 已吊销")
			} else if len(args) > 0 {
				auth.RevokeToken(args[0])
				fmt.Printf("Token %s... 已吊销\n", args[0][:16])
			} else {
				fmt.Fprintln(os.Stderr, "请指定 Token 或使用 --all 吊销所有")
				os.Exit(1)
			}
		},
	}
)

func init() {
	// 服务器参数
	rootCmd.Flags().StringVarP(&port, "port", "p", "8080", "服务端口")
	rootCmd.Flags().StringVar(&superMode, "super-mode", "token", "Super Mode 认证类型: token 或 password")
	rootCmd.Flags().StringVar(&superPass, "super-password", "", "Super Mode 密码 (当 type=password 时)")
	rootCmd.Flags().IntVar(&tokenTTL, "token-ttl", 24, "Token 有效期 (小时)")

	// Token 子命令参数
	tokenGenCmd.Flags().IntVar(&tokenTTL, "ttl", 24, "Token 有效期 (小时)")
	tokenRevokeCmd.Flags().BoolP("all", "a", false, "吊销所有 Token")

	tokenCmd.AddCommand(tokenGenCmd)
	tokenCmd.AddCommand(tokenStatusCmd)
	tokenCmd.AddCommand(tokenRevokeCmd)

	rootCmd.AddCommand(tokenCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	// 初始化认证管理器
	authConfig := auth.Config{
		Mode:     superMode,
		Password: superPass,
		TTL:      tokenTTL,
	}
	if err := auth.Initialize(authConfig); err != nil {
		log.Fatalf("初始化认证失败: %v", err)
	}

	// 初始化 K8s 客户端
	k8sClient, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("初始化 K8s 客户端失败: %v", err)
	}
	log.Println("✓ K8s 客户端初始化成功")

	// 设置 Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// 注册路由
	api.SetupRoutes(r, k8sClient)

	// 启动服务器
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	log.Printf("🚀 Ingress Portal 启动成功: http://localhost:%s", port)
	if superMode == "token" {
		token, _ := auth.GetActiveToken()
		if token != "" {
			log.Printf("🔐 Super Mode Token: %s", token[:16]+"...")
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}
	log.Println("✓ 服务器已关闭")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
