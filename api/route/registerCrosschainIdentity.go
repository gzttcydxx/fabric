package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/sdk"
	"github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func HandleRegisterCrosschainIdentity(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var didDoc models.DIDDoc
		if err := c.Bind(&didDoc); err != nil {
			c.JSON(400, gin.H{
				"statusCode": 1,
				"message":    "Bad Request",
			})
			return
		}
		statusCode, message := sdk.RegisterCrosschainIdentity(contract, didDoc)
		if statusCode != 0 {
			c.JSON(500, gin.H{
				"statusCode": statusCode,
				"message":    message,
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": statusCode,
				"message":    message,
			})
		}
	}
}
