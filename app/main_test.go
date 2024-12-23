package main_test

// import (
// 	"sync"
// 	"testing"
// 	"time"

// 	ecies "github.com/ecies/go/v2"
// 	"github.com/google/uuid"
// 	app "github.com/gzttcydxx/app"
// 	"github.com/stretchr/testify/require"
// )

// func TestRegisterChainID(t *testing.T) {
// 	privateKey, err := ecies.GenerateKey()
// 	require.NoError(t, err)
// 	publicKey := privateKey.PublicKey

// 	start := time.Now()

// 	var wg sync.WaitGroup
// 	for i := 0; i < app.Total; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			app.RegisterChainID(uuid.New().String(), publicKey.Hex(true), privateKey.Hex())
// 		}()
// 	}

// 	wg.Wait()

// 	elapsed := time.Since(start)
// 	tps := float64(app.Total) / elapsed.Seconds()
// 	t.Logf("Transactions per second: %f", tps)
// }
