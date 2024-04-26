package models

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/url"
// 	"regexp"
// 	"strings"
// )

// func ParseDID(did string) (*DID, error) {
// 	const idchar = `a-zA-Z0-9-_\.`
// 	regex := fmt.Sprintf(`^did:[a-z0-9]+:(:+|[:%s]+)*[%%:%s]+[^:]$`, idchar, idchar)

// 	r, err := regexp.Compile(regex)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to compile regex=%s (this should not have happened!). %w", regex, err)
// 	}

// 	if !r.MatchString(did) {
// 		return nil, fmt.Errorf(
// 			"invalid did: %s. Make sure it conforms to the DID syntax: https://w3c.github.io/did-core/#did-syntax", did)
// 	}

// 	parts := strings.SplitN(did, ":", 3)

// 	return &DID{
// 		Scheme:           "did",
// 		Method:           parts[1],
// 		MethodSpecificID: parts[2],
// 	}, nil
// }

// func ParseDIDURL(didURL string) (*DIDURL, error) {
// 	split := strings.IndexAny(didURL, "?/#")

// 	didPart := didURL
// 	pathQueryFragment := ""

// 	if split != -1 {
// 		didPart = didURL[:split]
// 		pathQueryFragment = didURL[split:]
// 	}

// 	retDID, err := ParseDID(didPart)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if pathQueryFragment == "" {
// 		return &DIDURL{
// 			DID:     *retDID,
// 			Queries: map[string][]string{},
// 		}, nil
// 	}

// 	hasPath := pathQueryFragment[0] == '/'

// 	if !hasPath {
// 		pathQueryFragment = "/" + pathQueryFragment
// 	}

// 	urlParts, err := url.Parse(pathQueryFragment)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse path, query, and fragment components of DID URL: %w", err)
// 	}

// 	ret := &DIDURL{
// 		DID:      *retDID,
// 		Queries:  urlParts.Query(),
// 		Fragment: urlParts.Fragment,
// 	}

// 	if hasPath {
// 		ret.Path = urlParts.Path
// 	}

// 	return ret, nil
// }

// func ParseDocument(data []byte) (*DIDDoc, error) {
// 	var raw rawDIDDoc

// 	err := json.Unmarshal(data, &raw)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal did doc: %w", err)
// 	}

// 	doc := &DIDDoc{
// 		Context:              raw.Context,
// 		ID:                   raw.ID,
// 		AlsoKnownAs:          any2Array(raw.AlsoKnownAs),
// 		Authentication:       any2Verification(raw.Authentication),
// 		AssertionMethod:      any2Verification(raw.AssertionMethod),
// 		CapabilityDelegation: any2Verification(raw.CapabilityDelegation),
// 		CapabilityInvocation: any2Verification(raw.CapabilityInvocation),
// 		KeyAgreement:         any2Verification(raw.KeyAgreement),
// 		Created:              raw.Created,
// 		Updated:              raw.Updated,
// 		Proof:                any2Proof(raw.Proof),
// 	}

// 	for _, vm := range raw.VerificationMethod {
// 		verificationMethod, err := parseVerificationMethod(vm)
// 		if err != nil {
// 			return nil, err
// 		}

// 		doc.VerificationMethod = append(doc.VerificationMethod, *verificationMethod)
// 	}

// 	return doc, nil
// }

// func parseVerificationMethod(vm map[string]interface{}) (*VerificationMethod, error) {
// 	id := any2String(vm["id"])
// 	if id == "" {
// 		return nil, fmt.Errorf("id is required")
// 	}

// 	relativeURL := false

// 	if strings.HasPrefix(id, "#") {
// 		relativeURL = true
// 	}

// 	return &VerificationMethod{
// 		ID:                 id,
// 		Type:               any2String(vm["type"]),
// 		Controller:         any2String(vm["controller"]),
// 		Value:              any2Bytes(vm["value"]),
// 		relativeURL:        relativeURL,
// 		PublicKeyMultibase: any2String(vm["publicKeyMultibase"]),
// 	}, nil
// }
