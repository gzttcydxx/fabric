package main

import (
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
