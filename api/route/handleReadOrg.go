package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/sdk"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func HandleReadOrg(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		org, err := sdk.ReadOrg(contract, c.Param("did"))
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to read org: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    org,
			})
		}
	}
}
