package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/gateway"
	"github.com/gzttcydxx/api/route"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Fabric API
// @version 0.0.1
// @description  测试
// @BasePath /v1/
func main() {
	contract, closeFunc := gateway.CreateNewConnection()
	defer closeFunc()

	// 配置 REST API 服务器
	r := gin.Default()
	// 设置信任的代理
	// r.SetTrustedProxies([]string{"traefik"})

	// 使用 CORS 中间件处理跨域问题，配置 CORS 参数
	r.Use(cors.New(cors.Config{
		// 允许的源地址（CORS中的Access-Control-Allow-Origin）
		// AllowOrigins: []string{"https://foo.com"},
		// 允许的 HTTP 方法（CORS中的Access-Control-Allow-Methods）
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		// 允许的 HTTP 头部（CORS中的Access-Control-Allow-Headers）
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		// 暴露的 HTTP 头部（CORS中的Access-Control-Expose-Headers）
		ExposeHeaders: []string{"Content-Length"},
		// 是否允许携带身份凭证（CORS中的Access-Control-Allow-Credentials）
		AllowCredentials: true,
		// 允许源的自定义判断函数，返回true表示允许，false表示不允许
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 允许你的开发环境
				return true
			} else if strings.HasPrefix(origin, "http://127.0.0.1") {
				return true
			}
			// 允许包含 "yourcompany.com" 的源
			return strings.Contains(origin, "a.gzttc.top")
		},
		// 用于缓存预检请求结果的最大时间（CORS中的Access-Control-Max-Age）
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/", route.HandleRoot)

	v1 := r.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.GET(":did", route.HandleReadUser(contract))
			users.GET("", route.HandleReadUsers(contract))
			users.POST("", route.HandleCreateUser(contract))
			users.PUT(":did", route.HandleUpdateUsers(contract))
			users.PATCH(":did", route.HandleUpdateUser(contract))
			users.DELETE(":did", route.HandleDeleteUser(contract))
		}

		orgs := v1.Group("/orgs")
		{
			orgs.GET(":did", route.HandleReadOrg(contract))
			orgs.GET("", route.HandleReadOrgs(contract))
			orgs.POST("", route.HandleCreateOrg(contract))
			orgs.PUT(":did", route.HandleUpdateOrgs(contract))
			orgs.PATCH(":did", route.HandleUpdateOrg(contract))
			orgs.DELETE(":did", route.HandleDeleteOrg(contract))
		}

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":80")
}
