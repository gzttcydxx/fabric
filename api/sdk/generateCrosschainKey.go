package sdk

import (
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func GenerateCrosschainKey(contract *client.Contract, did string) (int, string) {
	result, err := contract.EvaluateTransaction("GenerateCrosschainKey", did)
	if err != nil {
		return 1, fmt.Sprintf("failed to evaluate transaction: %v", err)
	}

	return 0, string(result)
}
