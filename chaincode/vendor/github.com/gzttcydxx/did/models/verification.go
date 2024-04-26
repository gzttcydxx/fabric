package models

import (
	"crypto/ed25519"
	"fmt"
	"io"

	"github.com/btcsuite/btcutil/base58"
)

type VerificationMethod struct {
	ID               DID    `json:"id,omitempty"`
	Type             string `json:"type,omitempty"`
	Controller       DID    `json:"controller,omitempty"`
	PublicKeyBase58  string `json:"publicKeyBase58,omitempty"`
	PrivateKeyBase58 string `json:"privateKeyBase58,omitempty"`
	// Value            []byte `json:"value,omitempty"`
}

func (v *VerificationMethod) GenerateKey(r io.Reader) error {
	publicKey, privateKey, err := ed25519.GenerateKey(r)
	if err != nil {
		return fmt.Errorf("failed to generate key: %v", err)
	}
	v.Type = "Ed25519VerificationKey2018"
	v.PublicKeyBase58 = base58.Encode(publicKey)
	v.PrivateKeyBase58 = base58.Encode(privateKey)

	return nil
}

func NewVerificationMethod(controller DID, publicKey, privateKey string) (*VerificationMethod, error) {
	vm := &VerificationMethod{
		Controller:       controller,
		Type:             "Ed25519VerificationKey2018",
		PublicKeyBase58:  publicKey,
		PrivateKeyBase58: privateKey,
	}
	vmID, err := NewDID(controller.ToString() + "#" + vm.PublicKeyBase58)
	if err != nil {
		return nil, err
	}
	vm.ID = *vmID
	return vm, nil
}

// type VerificationRelationship int

// const (
// 	VerificationRelationshipGeneral VerificationRelationship = iota

// 	Authentication

// 	AssertionMethod

// 	CapabilityDelegation

// 	CapabilityInvocation

// 	KeyAgreement
// )

// type Verification struct {
// 	VerificationMethod VerificationMethod       `json:"verificationMethod,omitempty"`
// 	Relationship       VerificationRelationship `json:"relationship,omitempty"`
// 	Embedded           bool                     `json:"embedded,omitempty"`
// }
