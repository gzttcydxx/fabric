package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"

// 	"github.com/gzttcydxx/app/gateway"
// )

// func main() {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	_, network, closeFunc := gateway.CreateNewConnection()
// 	defer closeFunc()

// 	signalChan := make(chan os.Signal, 1)
// 	signal.Notify(signalChan, os.Interrupt)

// 	// 开始监听链码事件
// 	events, err := network.ChaincodeEvents(ctx, "mychannel")
// 	if err != nil {
// 		fmt.Errorf("Failed to start chaincode event listening: %v", err)
// 	}

// 	fmt.Println("Started chaincode event listening...")

// 	// 监听事件
// 	go func() {
// 		for {
// 			select {
// 			case event, ok := <-events:
// 				if !ok {
// 					return
// 				}
// 				// handleEvent(event)
// 				fmt.Println(event.EventName)
// 			case <-signalChan:
// 				log.Println("Shutting down...")
// 				cancel()
// 				return
// 			}
// 		}
// 	}()

// 	<-ctx.Done()
// }
