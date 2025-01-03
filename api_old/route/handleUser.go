package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/sdk"
	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func HandleReadUser(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := sdk.ReadUser(contract, c.Param("did"))
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to read user: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    user,
			})
		}
	}
}

func HandleReadUsers(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := sdk.ReadUsers(contract)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to read users: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    users,
			})
		}
	}
}

func HandleCreateUser(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.Bind(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		result, err := sdk.CreateUser(contract, user)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to create user: " + err.Error(),
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

func HandleUpdateUser(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.Bind(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		err = sdk.UpdateUser(contract, c.Param("did"), user)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to update user: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "update success",
			})
		}
	}
}

func HandleUpdateUsers(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		err := c.Bind(&user)
		if err != nil {
			c.JSON(400, gin.H{
				"statusCode": 400,
				"message":    err,
			})
			return
		}
		err = sdk.UpdateUser(contract, c.Param("did"), user)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to update user: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "update success",
			})
		}
	}
}

func HandleDeleteUser(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := sdk.DeleteUser(contract, c.Param("did"))
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 500,
				"message":    "failed to delete user: " + err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"statusCode": 200,
				"message":    "delete success",
			})
		}
	}
}
