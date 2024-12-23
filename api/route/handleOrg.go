package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/sdk"
	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// @Summary     读取组织信息
// @Description 读取组织信息
// @Tags        org
// @Accept      json
// @Produce     json
// @Param       did path string true "组织 DID"
// Success 200 {object} models.Org "{"statusCode":200,"message":{}}"
// @Router /orgs [get]
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

func HandleCreateOrg(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var org models.Org
		err := c.Bind(&org)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		result, err := sdk.CreateOrg(contract, org)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to create org: " + err.Error(),
			})
			return
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    result,
			})
		}
	}
}

func HandleUpdateOrg(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var org models.Org
		err := c.Bind(&org)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		err = sdk.UpdateOrg(contract, c.Param("did"), org)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to update org: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "update success",
			})
		}
	}
}

func HandleUpdateOrgs(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var org models.Org
		err := c.Bind(&org)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		err = sdk.UpdateOrg(contract, c.Param("did"), org)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to update org: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "update success",
			})
		}
	}
}

func HandleDeleteOrg(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := sdk.DeleteOrg(contract, c.Param("did"))
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to delete org: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "delete success",
			})
		}
	}
}
