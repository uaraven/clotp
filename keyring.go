package main

import (
	"encoding/base64"
	"fmt"

	"github.com/99designs/keyring"
	"github.com/uaraven/gotp"
)

type Keys interface {
	ListOTPs() ([]OTPKey, error)
	AddKey(id string, name string, otpUri string) error
	RemoveById(id string) error
	RemoveByName(name string) error
}

type KeyringKeys struct {
	Keys
	ring keyring.Keyring
}

type OTPKey struct {
	gotp.OTPKeyData
	Id string
}

func NewKeys() (Keys, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "clotp",
	})
	if err != nil {
		return nil, err
	}
	return &KeyringKeys{
		ring: ring,
	}, nil
}

func (k *KeyringKeys) ListOTPs() ([]OTPKey, error) {
	result := make([]OTPKey, 0)
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
		otpKey := OTPKey{
			OTPKeyData: *otp,
			Id:         key,
		}
		result = append(result, otpKey)
	}
	return result, err
}

func (k *KeyringKeys) AddKey(id string, name string, otpUri string) error {
	base64uri := base64.StdEncoding.EncodeToString([]byte(otpUri))
	kitem := keyring.Item{
		Key:   id,
		Data:  []byte(base64uri),
		Label: name,
	}
	return k.ring.Set(kitem)
}

func (k *KeyringKeys) RemoveById(id string) error {
	return k.ring.Remove(id)
}

func (k *KeyringKeys) RemoveByName(name string) error {
	keys, err := k.ring.Keys()
	if err != nil {
		return err
	}
	for _, key := range keys {
		item, err := k.ring.Get(key)
		if err != nil {
			return err
		}
		if item.Label == name {
			return k.ring.Remove(item.Key)
		}
	}
	return fmt.Errorf("cannot find key with name=%s", name)
}
