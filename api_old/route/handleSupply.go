package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/sdk/services"
	"github.com/gzttcydxx/fabric/chaincode/models"
)

func RegisterSupplyRoutes(rg *gin.RouterGroup, service *services.SupplyService) {
	rg.POST("/available-stock", HandleGetAvailableStock(service))
	rg.POST("/accept-transaction", HandleAcceptTransaction(service))
}

func HandleGetAvailableStock(service *services.SupplyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var itemDemand models.ItemDemand
		if err := c.ShouldBindJSON(&itemDemand); err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    "invalid request body: " + err.Error(),
			})
			return
		}

		stock, err := service.GetAvailableStock(&itemDemand)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to get available stock: " + err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"statusCode": 200,
			"message":    stock,
		})
	}
}

func HandleAcceptTransaction(service *services.SupplyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Accept      bool               `json:"accept"`
			Transaction models.Transaction `json:"transaction"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    "invalid request body: " + err.Error(),
			})
			return
		}

		if err := service.AcceptTransaction(req.Accept, &req.Transaction); err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to accept/reject transaction: " + err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"statusCode": 200,
			"message":    "transaction processed successfully",
		})
	}
}
