package models

import (
	ecies "github.com/ecies/go/v2"
)

type Key struct {
	ID         DID    `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	Controller DID    `json:"controller,omitempty"`
	PublicKey  string `json:"publicKey,omitempty"`
	PrivateKey string `json:"privateKey,omitempty"`
}

func (k *Key) GenerateKey() error {
	privateKey, err := ecies.GenerateKey()
	if err != nil {
		return err
	}
	publicKey := privateKey.PublicKey

	k.PublicKey = publicKey.Hex(true)
	k.PrivateKey = privateKey.Hex()

	return nil
}

func NewKey(controller DID, publicKey, privateKey string) (*Key, error) {
	k := &Key{
		Controller: controller,
		Type:       "Secp256K1",
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	kID, err := NewDID(controller.ToString() + "#" + k.PublicKey)
	if err != nil {
		return nil, err
	}
	k.ID = *kID
	return k, nil
}
