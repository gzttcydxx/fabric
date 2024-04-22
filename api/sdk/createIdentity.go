package sdk

import (
	"fmt"

	ecies "github.com/ecies/go/v2"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

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
