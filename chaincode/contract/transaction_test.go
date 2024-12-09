package contract_test

import (
	"strings"
	"testing"

	"github.com/gzttcydxx/chaincode/contract"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/stretchr/testify/assert"
)

// Create a mock transaction context
type mockTransactionContext struct {
	shimtest.MockStub
}

func (m *mockTransactionContext) GetStub() shim.ChaincodeStubInterface {
	return &m.MockStub
}

// Add this method to implement the full interface
func (m *mockTransactionContext) GetClientIdentity() cid.ClientIdentity {
	return nil // For testing purposes, returning nil is often sufficient
}

func TestDeleteAllTransactions(t *testing.T) {
	// Create mock stub
	stub := shimtest.NewMockStub("mockStub", nil)
	ctx := &mockTransactionContext{*stub}

	// Start mock transaction
	stub.MockTransactionStart("tx1")

	// Prepare test data
	testData := map[string][]byte{
		"did:transaction:1": []byte("value1"),
		"did:transaction:2": []byte("value2"),
		"key3":              []byte("value3"),
	}

	// Write test data
	for k, v := range testData {
		err := stub.PutState(k, v)
		assert.NoError(t, err)
	}

	stub.MockTransactionEnd("tx1")

	// Start new transaction for delete operation
	stub.MockTransactionStart("tx2")

	// Create contract instance
	contract := new(contract.SmartContract)

	// Execute delete operation with mock context
	err := contract.DeleteAllTransactions(ctx)
	assert.NoError(t, err)

	// Verify all data has been deleted
	for k := range testData {
		value, err := stub.GetState(k)
		assert.NoError(t, err)

		if strings.HasPrefix(k, "did:transaction:") {
			assert.Nil(t, value, "Transaction data should be deleted: %s", k)
		} else {
			assert.NotNil(t, value, "Non-transaction data should not be deleted: %s", k)
		}
	}

	stub.MockTransactionEnd("tx2")
}
