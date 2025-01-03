package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/sdk/services"
	"github.com/gzttcydxx/fabric/chaincode/models"
)

// func RegisterDemandRoutes(rg *gin.RouterGroup, service *services.DemandService) {
// 	rg.POST("/init", HandleInitTransaction(service))
// 	rg.POST("/deal", HandleDealTransaction(service))
// 	rg.POST("/confirm", HandleConfirmTransaction(service))
// }

func HandleInitTransaction(service *services.DemandService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			DemandItemDid string `json:"demandItemDid"`
			DemandNum     string `json:"demandNum"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		demand, err := service.InitTransaction(req.DemandItemDid, req.DemandNum)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, demand)
	}
}

func HandleDealTransaction(service *services.DemandService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			SelfDid   string       `json:"selfDid"`
			Stock     models.Stock `json:"stock"`
			DemandNum string       `json:"demandNum"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		transaction, err := service.DealTransaction(req.SelfDid, req.Stock, req.DemandNum)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, transaction)
	}
}

func HandleConfirmTransaction(service *services.DemandService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Accept      bool               `json:"accept"`
			Transaction models.Transaction `json:"transaction"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := service.ConfirmTransaction(req.Accept, &req.Transaction); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "transaction confirmed"})
	}
}
