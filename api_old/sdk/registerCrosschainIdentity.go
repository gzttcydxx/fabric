package sdk

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func RegisterCrosschainIdentity(contract *client.Contract, didDoc models.DIDDoc) (int, string) {
	chainID := didDoc.ID.ChainID
	statusCode, result := ReadIdentity(contract, chainID)
	if statusCode != 0 {
		return 1, result
	}

	didDocJSON, err := json.Marshal(didDoc)
	if err != nil {
		return 1, fmt.Sprintf("failed to marshal didDoc: %v", err)
	}
	_, err = contract.SubmitTransaction("RegisterCrosschainIdentity", string(didDocJSON))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return 0, "The chainID " + chainID + " already registers"
		} else {
			return 1, fmt.Sprintf("failed to submit transaction: %v", err)
		}
	}

	return 0, "Register " + chainID + " Success"
}
