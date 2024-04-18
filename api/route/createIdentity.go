package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/gateway"
	"github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type createIdentityJson struct {
	Did string `json:"did"`
}

func HandleCreateIdentity(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json createIdentityJson
		if err := c.Bind(&json); err != nil {
			c.JSON(400, gin.H{
				"statusCode": 1,
				"message":    "Bad Request",
			})
			return
		}
		statusCode, message := gateway.CreateIdentity(contract, json.Did)
		var did models.DID
		did.FromString(json.Did)
		if statusCode != 0 {
			c.JSON(500, gin.H{
				"statusCode": statusCode,
				"message":    message,
				"did":        did,
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": statusCode,
				"message":    message,
				"did":        did,
			})
		}
	}
}
