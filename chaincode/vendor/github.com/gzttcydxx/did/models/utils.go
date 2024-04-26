package models

// import "time"

// func any2Verification(i interface{}) []Verification {
// 	if i == nil {
// 		return nil
// 	}

// 	is, ok := i.([]interface{})
// 	if !ok {
// 		return nil
// 	}

// 	var result []Verification

// 	for _, e := range is {
// 		if e != nil {
// 			result = append(result, any2VerificationElement(e))
// 		}
// 	}

// 	return result
// }

// func any2VerificationElement(i interface{}) Verification {
// 	switch e := i.(type) {
// 	case map[string]interface{}:
// 		vm, err := parseVerificationMethod(e)
// 		if err != nil {
// 			return Verification{}
// 		}

// 		return Verification{
// 			VerificationMethod: *vm,
// 		}
// 	case string:
// 		return Verification{
// 			VerificationMethod: VerificationMethod{
// 				ID: e,
// 			},
// 		}
// 	default:
// 		return Verification{}
// 	}
// }

// func any2Proof(i interface{}) []Proof {
// 	if i == nil {
// 		return nil
// 	}

// 	is, ok := i.([]interface{})
// 	if !ok {
// 		return nil
// 	}

// 	var result []Proof

// 	for _, e := range is {
// 		if e != nil {
// 			result = append(result, any2ProofElement(e))
// 		}
// 	}

// 	return result
// }

// func any2ProofElement(i interface{}) Proof {
// 	switch e := i.(type) {
// 	case map[string]interface{}:
// 		return Proof{
// 			Type:         any2String(e["type"]),
// 			Created:      any2Time(e["created"]),
// 			Creator:      any2String(e["creator"]),
// 			ProofValue:   any2Bytes(e["proofValue"]),
// 			Domain:       any2String(e["domain"]),
// 			Nonce:        any2Bytes(e["nonce"]),
// 			ProofPurpose: any2String(e["proofPurpose"]),
// 			relativeURL:  false,
// 		}
// 	default:
// 		return Proof{}
// 	}
// }

// func any2Time(i interface{}) *time.Time {
// 	if i == nil {
// 		return nil
// 	}

// 	switch e := i.(type) {
// 	case string:
// 		t, err := time.Parse(time.RFC3339, e)
// 		if err != nil {
// 			return nil
// 		}

// 		return &t
// 	default:
// 		return nil
// 	}
// }

// func any2Bytes(i interface{}) []byte {
// 	if i == nil {
// 		return nil
// 	}

// 	switch e := i.(type) {
// 	case string:
// 		return []byte(e)
// 	default:
// 		return nil
// 	}
// }

// func any2String(i interface{}) string {
// 	if i == nil {
// 		return ""
// 	}

// 	if e, ok := i.(string); ok {
// 		return e
// 	}

// 	return ""
// }

// func any2Array(i interface{}) []string {
// 	if i == nil {
// 		return nil
// 	}

// 	is, ok := i.([]interface{})
// 	if !ok {
// 		return nil
// 	}

// 	var result []string

// 	for _, e := range is {
// 		if e != nil {
// 			result = append(result, any2String(e))
// 		}
// 	}

// 	return result
// }
