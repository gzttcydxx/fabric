package route

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/gateway"
	"github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func HandleCreateCrosschainIdentityServer(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var didDoc models.DIDDoc
		if err := c.Bind(&didDoc); err != nil {
			c.JSON(400, gin.H{
				"statusCode": 1,
				"message":    "Bad Request",
			})
			return
		}
		statusCode, message := gateway.CreateCrosschainIdentityServer(contract, didDoc)
		if statusCode != 0 {
			c.JSON(500, gin.H{
				"statusCode": statusCode,
				"message":    message,
			})
		} else {
			var newdidDoc *models.DIDDoc
			err := json.Unmarshal([]byte(message), &newdidDoc)
			if err != nil {
				c.JSON(500, gin.H{
					"statusCode": 1,
					"message":    "failed to unmarshal did doc",
				})
			} else {
				c.JSON(200, gin.H{
					"statusCode": statusCode,
					"message":    newdidDoc,
				})
			}
		}
	}
}
