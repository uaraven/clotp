package cli

import (
	"crypto"
	"fmt"
	"strings"
)

func nameToHash(name string) (crypto.Hash, error) {
	ln := strings.ToLower(name)
	if ln == "sha-512" {
		return crypto.SHA512, nil
	} else if ln == "sha-256" {
		return crypto.SHA256, nil
	} else if ln == "sha-1" {
		return crypto.SHA1, nil
	} else {
		return crypto.SHA1, fmt.Errorf("unsupported hash: %s", name)
	}
}
