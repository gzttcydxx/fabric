package models

type Context []string

type DIDDoc struct {
	Context     Context `json:"@context"`
	ID          DID     `json:"id"`
	AlsoKnownAs DID     `json:"alsoKnownAs"`
	Key         []Key
	// VerificationMethod []VerificationMethod
	// Created              *time.Time
	// Updated              *time.Time
	// Proof                []Proof
}

func NewDIDDoc(did DID, publicKey, privateKey string) (*DIDDoc, error) {
	diddoc := &DIDDoc{
		Context: []string{"https://www.w3.org/ns/did/v1", "https://w3id.org/security/suites/ed25519-2018/v1"},
		ID:      did,
	}

	// vm, err := NewVerificationMethod(did, publicKey, privateKey)
	// if err != nil {
	// 	return nil, err
	// }
	// diddoc.VerificationMethod = append(diddoc.VerificationMethod, *vm)

	k, err := NewKey(did, publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	diddoc.Key = append(diddoc.Key, *k)

	return diddoc, nil
}

// type rawDIDDoc struct {
// 	Context              Context                  `json:"@context,omitempty"`
// 	ID                   string                   `json:"id,omitempty"`
// 	AlsoKnownAs          []interface{}            `json:"alsoKnownAs,omitempty"`
// 	VerificationMethod   []map[string]interface{} `json:"verificationMethod,omitempty"`
// 	Authentication       []interface{}            `json:"authentication,omitempty"`
// 	AssertionMethod      []interface{}            `json:"assertionMethod,omitempty"`
// 	CapabilityDelegation []interface{}            `json:"capabilityDelegation,omitempty"`
// 	CapabilityInvocation []interface{}            `json:"capabilityInvocation,omitempty"`
// 	KeyAgreement         []interface{}            `json:"keyAgreement,omitempty"`
// 	Created              *time.Time               `json:"created,omitempty"`
// 	Updated              *time.Time               `json:"updated,omitempty"`
// 	Proof                []interface{}            `json:"proof,omitempty"`
// }

// // UnmarshalJSON unmarshals a DID Document.
// func (doc *DIDDoc) UnmarshalJSON(data []byte) error {
// 	_doc, err := ParseDocument(data)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse did doc: %w", err)
// 	}

// 	*doc = *_doc

// 	return nil
// }
