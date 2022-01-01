package main

import (
	"encoding/base64"

	"github.com/99designs/keyring"
	"github.com/uaraven/gotp"
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

func (k *Keys) ListOTPs() ([]gotp.OTPKeyData, error) {
	result := make([]gotp.OTPKeyData, 0)
	keys, err := k.ring.Keys()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		item, err := k.ring.Get(key)
		if err != nil {
			return nil, err
		}
		uriBytes, err := base64.StdEncoding.DecodeString(string(item.Data))
		if err != nil {
			return nil, err
		}
		otp, err := gotp.OTPFromUri(string(uriBytes))
		if err != nil {
			return nil, err
		}
		result = append(result, *otp)
	}
	return result, err
}

func (k *Keys) AddKey(name string, otpUri string) error {
	base64uri := base64.StdEncoding.EncodeToString([]byte(otpUri))
	kitem := keyring.Item{
		Key:   name,
		Data:  []byte(base64uri),
		Label: name,
	}
	return k.ring.Set(kitem)
}
