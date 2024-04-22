package gateway

import (
	"encoding/json"
	"fmt"
	"strings"

	// "github.com/btcsuite/btcutil/base58"
	ecies "github.com/ecies/go/v2"
	"github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// func getAllAssets(contract *client.Contract) {
// 	fmt.Println("Evaluate Transaction: GetAllAssets, function returns all the current assets on the ledger")

// 	evaluateResult, err := contract.EvaluateTransaction("GetAllProjects")
// 	if err != nil {
// 		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
// 	}
// 	result := formatJSON(evaluateResult)

// 	fmt.Printf("*** Result:%s\n", result)
// }

func CreateIdentity(contract *client.Contract, did string) (int, string) {
	// 生成 Ed25519 密钥对
	// publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	// if err != nil {
	// 	return 1, fmt.Sprintf("failed to generate key pair: %v", err)
	// }
	// // 使用 Base58 编码
	// publicKeyBase58 := base58.Encode(publicKey)
	// privateKeyBase58 := base58.Encode(privateKey)
	privateKey, err := ecies.GenerateKey()
	if err != nil {
		return 1, fmt.Sprintf("failed to generate key pair: %v", err)
	}
	publicKey := privateKey.PublicKey

	_, err = contract.SubmitTransaction("CreateIdentity", did, publicKey.Hex(true), privateKey.Hex())
	if err != nil {
		return 1, fmt.Sprintf("failed to submit transaction: %v", err)
	}

	return 0, "Create " + did + " Success"
}

func ReadIdentity(contract *client.Contract, did string) (int, string) {
	result, err := contract.EvaluateTransaction("ReadIdentity", did)
	if err != nil {
		return 1, fmt.Sprintf("failed to evaluate transaction: %v", err)
	}

	return 0, string(result)
}

// TODO: 生成跨链身份
func CreateCrosschainIdentityServer(contract *client.Contract, didDoc models.DIDDoc) (int, string) {
	chainID := didDoc.ID.ChainID
	statusCode, result := ReadIdentity(contract, chainID)
	if statusCode != 0 {
		return 1, result
	}

	didDocJSON, err := json.Marshal(didDoc)
	if err != nil {
		return 1, fmt.Sprintf("failed to marshal didDoc: %v", err)
	}
	_, err = contract.SubmitTransaction("CreateCrosschainIdentityServer", string(didDocJSON))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return 0, "The chainID " + chainID + " already exists"
		} else {
			return 1, fmt.Sprintf("failed to submit transaction: %v", err)
		}
	}

	return 0, "Create " + chainID + " Success"
}

// func CreateCrosschainIdentityServer(contract *client.Contract, didDoc models.DIDDoc) (int, string) {
// 	privateKey, err := ecies.GenerateKey()
// 	if err != nil {
// 		return 1, fmt.Sprintf("failed to generate key pair: %v", err)
// 	}
// 	publicKey := privateKey.PublicKey

// 	didDocRaw, err := json.Marshal(didDoc)
// 	if err != nil {
// 		return 1, fmt.Sprintf("failed to marshal didDoc: %v", err)
// 	}

// 	chainID, err := uuid.NewRandom()
// 	if err != nil {
// 		return 1, fmt.Sprintf("failed to generate uuid: %v", err)
// 	}

// 	result, err := contract.SubmitTransaction("CreateCrosschainIdentityServer", string(didDocRaw), publicKey.Hex(true), privateKey.Hex(), chainID.String())
// 	if err != nil {
// 		return 1, fmt.Sprintf("failed to evaluate transaction: %v", err)
// 	}

// 	return 0, string(result)
// }

// TODO: 生成跨链会话密钥
func GenerateCrosschainKey(contract *client.Contract, did string) (int, string) {
	result, err := contract.EvaluateTransaction("GenerateCrosschainKey", did)
	if err != nil {
		return 1, fmt.Sprintf("failed to evaluate transaction: %v", err)
	}

	return 0, string(result)
}
