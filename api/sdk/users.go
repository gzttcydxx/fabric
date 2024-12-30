package sdk

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	modelDID "github.com/gzttcydxx/did/models"
	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func ReadUser(contract *client.Contract, did string) (models.User, error) {
	var user models.User
	result, err := contract.EvaluateTransaction("ReadUser", did)
	if err != nil {
		return user, fmt.Errorf("failed to evaluate transaction: %v", err)
	}

	err = json.Unmarshal(result, &user)
	if err != nil {
		return user, fmt.Errorf("failed to unmarshal did doc: %v", err)
	}

	return user, nil
}

func ReadUsers(contract *client.Contract) ([]models.User, error) {
	var users []models.User
	result, err := contract.EvaluateTransaction("ReadUsers")
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %v", err)
	}

	err = json.Unmarshal(result, &users)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user infos: %v", err)
	}

	return users, nil
}

func CreateUser(contract *client.Contract, user models.User) (string, error) {
	spid := sha256.Sum256([]byte(user.Name + user.Role + user.Org.Did.ToString()))
	did := "did:sha256:" + hex.EncodeToString(spid[:])
	didValue, _ := modelDID.NewDID(did)
	user.Did = *didValue
	userJson, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user: %v", err)
	}

	_, err = contract.SubmitTransaction("CreateUser", string(userJson))
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction: %v", err)
	}

	return did, nil
}

func UpdateUser(contract *client.Contract, did string, user models.User) error {
	user.Did.FromString(did)
	userJson, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	_, err = contract.SubmitTransaction("UpdateUser", string(userJson))
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}

	return nil
}

func DeleteUser(contract *client.Contract, did string) error {
	_, err := contract.SubmitTransaction("DeleteUser", did)
	if err != nil {
		return fmt.Errorf("failed to submit transaction: %v", err)
	}

	return nil
}
