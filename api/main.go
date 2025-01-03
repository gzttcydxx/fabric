package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hyperledger/fabric-gateway/pkg/client"

	"github.com/gzttcydxx/newapi/gateway"
	"github.com/gzttcydxx/newapi/models"
	"github.com/gzttcydxx/newapi/routes"
)

func registerRoutes(api huma.API, contract *client.Contract, version string) {
	// 注册组织的CRUD路由
	routes.RegisterCRUD[models.Org](api, contract, fmt.Sprintf("/%s/orgs", version), "organization", models.CRUDMethods{
		Create: "CreateOrg",
		Read:   "ReadOrg",
		Query:  "QueryOrgs",
		Update: "UpdateOrg",
		Delete: "DeleteOrg",
	})

	// 注册零件的CRUD路由
	routes.RegisterCRUD[models.Part](api, contract, fmt.Sprintf("/%s/parts", version), "part", models.CRUDMethods{
		Create: "CreatePart",
		Read:   "ReadPart",
		Query:  "QueryParts",
		Update: "UpdatePart",
		Delete: "DeletePart",
	})

	// 注册产品的CRUD路由
	routes.RegisterCRUD[models.Product](api, contract, fmt.Sprintf("/%s/products", version), "product", models.CRUDMethods{
		Create: "CreateProduct",
		Read:   "ReadProduct",
		Query:  "QueryProducts",
		Update: "UpdateProduct",
		Delete: "DeleteProduct",
	})
}

func main() {
	version := "v1"

	contract, closeFunc := gateway.CreateNewConnection()
	defer closeFunc()

	router := chi.NewMux()

	// 添加请求日志中间件
	router.Use(middleware.Logger)

	api := humachi.New(router, huma.DefaultConfig("Fabric Transaction API", "1.0.0"))

	registerRoutes(api, contract, version)

	// 记录服务器启动信息
	log.Printf("Server starting on :80")
	if err := http.ListenAndServe("0.0.0.0:80", router); err != nil {
		log.Fatal(err)
	}
}
