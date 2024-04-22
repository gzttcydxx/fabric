package route

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gzttcydxx/did/models"

	"github.com/gzttcydxx/api/sdk"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleGenerateCrosschainKey(contract *client.Contract) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 1,
				"message":    fmt.Sprintf("failed to upgrade connection %v", err),
			})
			return
		}
		defer conn.Close()

		statusCode, message := sdk.ReadIdentity(contract, "did:example:admin")
		if statusCode != 0 {
			c.JSON(500, gin.H{
				"statusCode": statusCode,
				"message":    message,
			})
			return
		}

		var DIDdoc *models.DIDDoc
		err = json.Unmarshal([]byte(message), &DIDdoc)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 1,
				"message":    fmt.Sprintf("failed to unmarshal did doc %v", err),
			})
			return
		}

		// sk := DIDdoc.VerificationMethod[0].PrivateKeyBase58

		var messageByte []byte

		// DID_{Alice}, DID_{Bob}, N, T
		_, messageByte, err = conn.ReadMessage()
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": 1,
				"message":    "failed to upgrade connection",
			})
			return
		}

		message = string(messageByte)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// fmt.Printf("recv: %s\n", message)
			statusCode, result := sdk.GenerateCrosschainKey(contract, string(message))
			if statusCode != 0 {
				conn.WriteMessage(websocket.TextMessage, []byte(result))
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte(result))
			}
		}
	}
}
