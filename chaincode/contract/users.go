package contract

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, did string) (models.User, error) {
	var user models.User
	result, err := ctx.GetStub().GetState(did)
	if err != nil {
		return user, fmt.Errorf("failed to get user info: %v", err)
	}
	if result == nil {
		return user, fmt.Errorf("the user %s does not exist", did)
	}

	err = json.Unmarshal(result, &user)
	if err != nil {
		return user, fmt.Errorf("failed to unmarshal user info: %v", err)
	}

	return user, nil
}

func (s *SmartContract) ReadUsers(ctx contractapi.TransactionContextInterface) ([]models.User, error) {
	queryString := `{"selector":{"role":{"$exists":true}}}`

	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer iterator.Close()

	var users []models.User
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var user models.User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, userJson string) error {
	var user models.User
	err := json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %v", err)
	}

	did := user.Did
	org := user.Org

	userResult, _ := s.ReadUser(ctx, did.ToString())
	if userResult.Did.SpecificID != "" {
		return fmt.Errorf("the user %s already exists", did)
	}
	// if err != nil {
	// 	return fmt.Errorf("failed to get user info: %v", err)
	// }

	orgResult, err := s.ReadOrg(ctx, org.Did.ToString())
	if orgResult.Did.SpecificID == "" {
		return fmt.Errorf("the org %s does not exist", org.Did)
	}
	if err != nil {
		return fmt.Errorf("failed to get org info: %v", err)
	}

	err = ctx.GetStub().PutState(did.ToString(), []byte(userJson))
	if err != nil {
		return fmt.Errorf("failed to put user: %v", err)
	}

	return nil
}

func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, userJson string) error {
	var user models.User
	err := json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %v", err)
	}

	did := user.Did

	readUser, err := s.ReadUser(ctx, did.ToString())
	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}
	if readUser.Did.SpecificID == "" {
		return fmt.Errorf("the user %s does not exist", did)
	}

	if user.Name != "" {
		readUser.Name = user.Name
	}
	if user.Role != "" {
		readUser.Role = user.Role
	}
	if user.Org.Did.SpecificID != "" {
		orgResult, err := s.ReadOrg(ctx, user.Org.Did.ToString())
		if err != nil {
			return fmt.Errorf("failed to get org info: %v", err)
		}
		if orgResult.Did.SpecificID == "" {
			return fmt.Errorf("the org %s does not exist", user.Org.Did)
		}
		readUser.Org = user.Org
	}

	readJson, err := json.Marshal(readUser)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	err = ctx.GetStub().PutState(did.ToString(), []byte(readJson))
	if err != nil {
		return fmt.Errorf("failed to put user: %v", err)
	}

	return nil
}

func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, did string) error {
	readUser, err := s.ReadUser(ctx, did)
	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}
	if readUser.Did.SpecificID == "" {
		return fmt.Errorf("the user %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}
