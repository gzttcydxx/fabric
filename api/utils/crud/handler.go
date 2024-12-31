package crud

import (
	"github.com/gin-gonic/gin"
)

// HandleRead 通用读取处理函数
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
