package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gzttcydxx/api/gateway"
	"github.com/gzttcydxx/api/route"
)

// func createIdentity(i int, j int, contract *client.Contract, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	did := fmt.Sprintf("did:example:%d_%d", i, j)
// 	status, _ := gateway.CreateIdentity(contract, did)
// 	if status != 0 {
// 		fmt.Printf("did:example:%d_%d CreateIdentity failed!\n", i, j)
// 	}
// }

func main() {
	contract, closeFunc := gateway.CreateNewConnection()
	defer closeFunc()

	// var wg sync.WaitGroup
	// sem := make(chan struct{}, 20)
	// for i := 0; i < 100; i++ {
	// 	for j := 0; j < 100; j++ {
	// 		wg.Add(1)
	// 		sem <- struct{}{}
	// 		go func(i, j int) {
	// 			createIdentity(i, j, contract, &wg)
	// 			<-sem
	// 		}(i, j)
	// 	}
	// }
	// wg.Wait()
	// for len(sem) > 0 {
	// 	time.Sleep(1 * time.Second)
	// }

	// 配置 REST API 服务器
	r := gin.Default()
	// 设置信任的代理
	// r.SetTrustedProxies([]string{"traefik"})

	r.GET("/", route.HandleRoot)
	r.POST("/create_identity", route.HandleCreateIdentity(contract))
	r.GET("/read_identity", route.HandleReadIdentity(contract))
	r.POST("/create_crosschain_identity", route.HandleRegisterCrosschainIdentity(contract))
	r.GET("/generate_crosschain_key", route.HandleGenerateCrosschainKey(contract))

	r.Run(":80")

	// ciphertext, err := ecies.Encrypt(k.PublicKey, []byte("THIS IS THE TEST"))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("plaintext encrypted: %x\n", ciphertext)

	// plaintext, err := ecies.Decrypt(k, ciphertext)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("ciphertext decrypted: %s\n", string(plaintext))
}
