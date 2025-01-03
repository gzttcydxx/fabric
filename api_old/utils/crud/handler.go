package crud

import (
	"github.com/gin-gonic/gin"
)

// HandleRead 通用读取处理函数
// @Summary 获取实体
// @Description 根据 DID 获取实体信息
// @Tags crud
// @Accept json
// @Produce json
// @Param type path string true "实体类型" Enums(orgs,parts,products)
// @Param did path string true "实体 DID"
// @Success 200 {object} models.Response{message=models.Org} "组织数据" SchemaExample(orgs)
// @Success 200 {object} models.Response{message=models.Part} "零件数据" SchemaExample(parts)
// @Success 200 {object} models.Response{message=models.Product} "产品数据" SchemaExample(products)
// @Failure 404 {object} models.Response
// @Router /v1/{type}/{did} [get]
func HandleRead[T any](service *CRUDService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		entity, err := service.Read(c.Param("did"))
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to read: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    entity,
			})
		}
	}
}

// HandleQuery 通用查询处理函数
// @Summary 查询实体列表
// @Description 根据条件查询实体列表
// @Tags crud
// @Accept json
// @Produce json
// @Param type path string true "实体类型" Enums(orgs,parts,products)
// @Param org body models.Org false "组织查询条件" SchemaExample(orgs)
// @Param part body models.Part false "零件查询条件" SchemaExample(parts)
// @Param product body models.Product false "产品查询条件" SchemaExample(products)
// @Success 200 {array} models.Response{message=models.Org} "组织列表" SchemaExample(orgs)
// @Success 200 {array} models.Response{message=models.Part} "零件列表" SchemaExample(parts)
// @Success 200 {array} models.Response{message=models.Product} "产品列表" SchemaExample(products)
// @Failure 400 {object} models.Response
// @Router /{type}/query [post]
func HandleQuery[T any](service *CRUDService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建新的实例
		var params T

		// 解析 JSON 请求体
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    "invalid request body: " + err.Error(),
			})
			return
		}

		// 调用查询函数
		results, err := service.Query(params)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to query: " + err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"statusCode": 200,
			"message":    results,
		})
	}
}

// HandleCreate 通用创建处理函数
// @Summary 创建实体
// @Description 创建新的实体
// @Tags crud
// @Accept json
// @Produce json
// @Param type path string true "实体类型" Enums(orgs,parts,products)
// @Param org body models.Org true "创建组织" SchemaExample(orgs)
// @Param part body models.Part true "创建零件" SchemaExample(parts)
// @Param product body models.Product true "创建产品" SchemaExample(products)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /{type} [post]
func HandleCreate[T any](service *CRUDService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params T
		err := c.Bind(&params)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		err = service.Create(params)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to create: " + err.Error(),
			})
			return
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "created successfully",
			})
		}
	}
}

// HandleUpdate 通用更新处理函数
// @Summary 更新实体
// @Description 更新实体信息
// @Tags crud
// @Accept json
// @Produce json
// @Param type path string true "实体类型" Enums(orgs,parts,products)
// @Param org body models.Org false "组织更新数据" SchemaExample(orgs)
// @Param part body models.Part false "零件更新数据" SchemaExample(parts)
// @Param product body models.Product false "产品更新数据" SchemaExample(products)
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /{type} [patch]
func HandleUpdate[T any](service *CRUDService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params T
		err := c.Bind(&params)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		err = service.Update(params)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to update: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "updated successfully",
			})
		}
	}
}

// HandleDelete 通用删除处理函数
// @Summary 删除实体
// @Description 删除指定实体
// @Tags crud
// @Accept json
// @Produce json
// @Param type path string true "实体类型" Enums(orgs,parts,products)
// @Param did path string true "实体 DID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /{type}/{did} [delete]
func HandleDelete[T any](service *CRUDService[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := service.Delete(c.Param("did"))
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to delete: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "deleted successfully",
			})
		}
	}
}
