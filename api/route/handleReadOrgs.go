package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/sdk"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func HandleReadOrgs(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgs, err := sdk.ReadOrgs(contract)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to read orgs: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    orgs,
			})
		}
	}
}
