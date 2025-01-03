package route

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/gzttcydxx/api/sdk"
// 	"github.com/hyperledger/fabric-gateway/pkg/client"
// )

// func HandleInitTransaction(contract *client.Contract) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tx, err := sdk.InitTransaction(contract, c.Param("demandItemDid"), c.Param("demandNum"))
// 		if err != nil {
// 			c.JSON(500, gin.H{
// 				"statusCode": 500,
// 				"message":    "failed to init transaction: " + err.Error(),
// 			})
// 		} else {
// 			c.JSON(200, gin.H{
// 				"statusCode": 200,
// 				"message":    tx,
// 			})
// 		}
// 	}
// }
