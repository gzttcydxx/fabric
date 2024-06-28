package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/models"
	"github.com/gzttcydxx/api/sdk"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func HandleCreateOrg(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var org models.Org
		err := c.Bind(&org)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		result, err := sdk.CreateOrg(contract, org)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to create org: " + err.Error(),
			})
			return
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    result,
			})
		}
	}
}
