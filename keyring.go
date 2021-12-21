package main

import (
	"github.com/99designs/keyring"
)

type Keys struct {
	ring keyring.Keyring
}

func NewKeys() (*Keys, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "clotp",
	})
	if err != nil {
		return nil, err
	}
	return &Keys{
		ring: ring,
	}, nil
}

func (k *Keys) ListKeys() ([]string, error) {
	return k.ring.Keys()
}
