package route

import "github.com/gin-gonic/gin"

func HandleRoot(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}
