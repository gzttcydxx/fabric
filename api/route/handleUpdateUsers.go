package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/models"
	"github.com/gzttcydxx/api/sdk"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func HandleUpdateUsers(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.Bind(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		err = sdk.UpdateUser(contract, c.Param("did"), user)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to update user: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "update success",
			})
		}
	}
}
