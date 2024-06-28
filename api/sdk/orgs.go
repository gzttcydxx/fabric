package sdk

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/api/models"
	modelDID "github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func ReadOrg(contract *client.Contract, did string) (models.Org, error) {
	var org models.Org
	result, err := contract.EvaluateTransaction("ReadOrg", did)
	if err != nil {
		return org, fmt.Errorf("failed to evaluate transaction: %v", err)
	}

	err = json.Unmarshal(result, &org)
	if err != nil {
		return org, fmt.Errorf("failed to unmarshal did doc: %v", err)
	}

	return org, nil
}

func ReadOrgs(contract *client.Contract) ([]models.Org, error) {
	var orgs []models.Org
	result, err := contract.EvaluateTransaction("ReadOrgs")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %v", err)
	}

	err = json.Unmarshal(result, &orgs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal org infos: %v", err)
	}

	return orgs, nil
}

func CreateOrg(contract *client.Contract, org models.Org) (string, error) {
	spid := sha256.Sum256([]byte(org.Name))
	did := "did:sha256:" + hex.EncodeToString(spid[:])
	didValue, _ := modelDID.NewDID(did)
	org.Did = *didValue
	orgJson, err := json.Marshal(org)
	if err != nil {
		return "", fmt.Errorf("failed to marshal org: %v", err)
	}

	_, err = contract.SubmitTransaction("CreateOrg", string(orgJson))
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction: %v", err)
	}

	return did, nil
}

func UpdateOrg(contract *client.Contract, did string, org models.Org) error {
	org.Did.FromString(did)
	orgJson, err := json.Marshal(org)
	if err != nil {
		return fmt.Errorf("failed to marshal org: %v", err)
	}

	_, err = contract.SubmitTransaction("UpdateOrg", string(orgJson))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}

	return nil
}

func DeleteOrg(contract *client.Contract, did string) error {
	_, err := contract.SubmitTransaction("DeleteOrg", did)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}

	return nil
}
