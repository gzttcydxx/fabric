package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gzttcydxx/app/gateway"
	"github.com/gzttcydxx/did/models"
)

type ResponseMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

var (
	Total            = 100
	ConcurrencyLimit = 10
	Semaphore        = make(chan struct{}, ConcurrencyLimit)
)

func registerChainID() {
	Semaphore <- struct{}{}
	defer func() { <-Semaphore }()

	contract, closeFunc := gateway.CreateNewConnection()
	defer closeFunc()

	chainIDbyte, err := contract.EvaluateTransaction("GetChainID")
	if err != nil {
		panic(err)
	}
	chainID := string(chainIDbyte)

	chainIDverifiedByte, err := contract.EvaluateTransaction("GetChainIDVerified")
	if err != nil {
		panic(err)
	}
	chainIDverified := string(chainIDverifiedByte)

	if chainIDverified != "true" {
		adminDIDDocJSON, err := contract.EvaluateTransaction("ReadIdentity", "did:example:admin")
		if err != nil {
			panic(err)
		}
		if adminDIDDocJSON == nil {
			panic("the identity did:example:admin does not exist")
		}

		var adminDIDDoc models.DIDDoc

		json.Unmarshal(adminDIDDocJSON, &adminDIDDoc)
		adminDIDDoc.ID.ChainID = string(chainID)
		adminDIDDocJSON, err = json.Marshal(adminDIDDoc)
		if err != nil {
			panic(err)
		}
		resp, err := http.Post("https://api.rc.gzttc.top/create_crosschain_identity", "application/json", bytes.NewBuffer(adminDIDDocJSON))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var responseData ResponseMessage
		jsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(jsonData, &responseData)
		if responseData.StatusCode != 0 {
			panic(fmt.Sprintf("failed to register chainID: %s", responseData.Message))
		}

		_, err = contract.SubmitTransaction("UpdateIdentity", string(adminDIDDocJSON))
		if err != nil {
			panic(err)
		}

		_, err = contract.SubmitTransaction("SetChainIDVerified", "true")
		if err != nil {
			panic(err)
		}
	}
}

func RegisterChainID(chainID, publicKey, privateKey string) {
	Semaphore <- struct{}{}
	defer func() { <-Semaphore }()

	// contract, closeFunc := gateway.CreateNewConnection()
	// defer closeFunc()

	adminDID, _ := models.NewDID("did:example:admin-" + chainID)
	adminDID.ChainID = chainID

	adminDIDDoc, _ := models.NewDIDDoc(*adminDID, publicKey, privateKey)

	adminDIDDocJSON, _ := json.Marshal(adminDIDDoc)

	resp, err := http.Post("https://api.rc.gzttc.top/create_crosschain_identity", "application/json", bytes.NewBuffer(adminDIDDocJSON))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var responseData ResponseMessage
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(jsonData, &responseData)
	if responseData.StatusCode != 0 {
		panic(fmt.Sprintf("failed to register chainID: %s", responseData.Message))
	}

	// _, err = contract.SubmitTransaction("UpdateIdentity", string(adminDIDDocJSON))
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = contract.SubmitTransaction("SetChainIDVerified", "true")
	// if err != nil {
	// 	panic(err)
	// }
}
