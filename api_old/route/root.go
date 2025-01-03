package route

import (
	"github.com/gin-gonic/gin"
)

// @Summary 健康检查接口
// @Description 返回服务状态信息，用于检查 API 服务是否正常运行
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{message=string} "服务正常"
// @Router / [get]
func HandleRoot(c *gin.Context) {
	c.JSON(200, gin.H{
		"statusCode": 200,
		"message":    "OK",
	})
}
