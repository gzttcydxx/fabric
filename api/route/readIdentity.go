package route

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/gateway"
	"github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func HandleReadIdentity(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode, message := gateway.ReadIdentity(contract, c.Query("did"))
		if statusCode != 0 {
			c.JSON(500, gin.H{
				"statusCode": statusCode,
				"message":    message,
			})
		} else {
			var DIDdoc *models.DIDDoc
			err := json.Unmarshal([]byte(message), &DIDdoc)
			if err != nil {
				c.JSON(500, gin.H{
					"statusCode": 1,
					"message":    "failed to unmarshal did doc",
				})
			} else {
				c.JSON(200, gin.H{
					"statusCode": statusCode,
					"message":    DIDdoc,
				})
			}
		}
	}
}
