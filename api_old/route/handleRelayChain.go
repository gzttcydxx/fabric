package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/sdk/services"
	"github.com/gzttcydxx/fabric/chaincode/models"
)

func RegisterRelayChainRoutes(rg *gin.RouterGroup, service *services.RelayChainService) {
	rg.POST("/filter-stock", HandleFilterDemandStock(service))
	rg.POST("/send-to-supply", HandleSend2SupplyTransaction(service))
	rg.POST("/send-to-demand", HandleSend2DemandTransaction(service))
}

func HandleFilterDemandStock(service *services.RelayChainService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var itemDemand models.ItemDemand
		if err := c.ShouldBindJSON(&itemDemand); err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    "invalid request body: " + err.Error(),
			})
			return
		}

		stock, err := service.FilterDemandStock(&itemDemand)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to filter stock: " + err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"statusCode": 200,
			"message":    stock,
		})
	}
}

func HandleSend2SupplyTransaction(service *services.RelayChainService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction models.Transaction
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    "invalid request body: " + err.Error(),
			})
			return
		}

		result, err := service.Send2SupplyTransaction(&transaction)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to send transaction: " + err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"statusCode": 200,
			"message":    result,
		})
	}
}

func HandleSend2DemandTransaction(service *services.RelayChainService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var transaction models.Transaction
		if err := c.ShouldBindJSON(&transaction); err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    "invalid request body: " + err.Error(),
			})
			return
		}

		result, err := service.Send2DemandTransaction(&transaction)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to send transaction: " + err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"statusCode": 200,
			"message":    result,
		})
	}
}
