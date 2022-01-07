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
	GetById(id string) (*OTPKey, error)
	GetByName(name string) (*OTPKey, error)
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

func (k *KeyringKeys) dataToOtpKey(item keyring.Item) (*OTPKey, error) {
	if len(item.Data) != 0 {
		uriBytes, err := base64.StdEncoding.DecodeString(string(item.Data))
		if err != nil {
			return nil, err
		}
		otp, err := gotp.OTPFromUri(string(uriBytes))
		if err != nil {
			return nil, err
		}
		return &OTPKey{OTPKeyData: *otp, Id: item.Key}, nil
	}
	return nil, nil
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
		otp, err := k.dataToOtpKey(item)
		if err != nil {
			return nil, err
		}
		if otp != nil {
			result = append(result, *otp)
		}
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

func (k *KeyringKeys) GetById(id string) (*OTPKey, error) {
	item, err := k.ring.Get(id)
	if err != nil {
		return nil, err
	}
	key, err := k.dataToOtpKey(item)
	if key == nil && err == nil {
		return nil, fmt.Errorf("cannot find code with id=%s", id)
	} else {
		return key, err
	}
}

func (k *KeyringKeys) GetByName(name string) (*OTPKey, error) {
	keys, err := k.ring.Keys()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		item, err := k.ring.Get(key)
		if err != nil {
			return nil, err
		}
		if item.Label == name {
			return k.dataToOtpKey(item)
		}
	}
	return nil, fmt.Errorf("cannot find code with name=%s", name)
}

func (k *KeyringKeys) RemoveById(id string) error {
	// return k.ring.Remove(id)
	// remove always fails, so we just clear the data for now
	item, err := k.ring.Get(id)
	if err != nil {
		return err
	}
	item.Data = make([]byte, 0)
	return k.ring.Set(item)
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
			// remove always fails, so we just clear the data for now
			// return k.ring.Remove(item.Key)
			item.Data = make([]byte, 0)
			return k.ring.Set(item)
		}
	}
	return fmt.Errorf("cannot find key with name=%s", name)
}
