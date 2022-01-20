package keyrings

import (
	"fmt"

	"github.com/99designs/keyring"
	"github.com/uaraven/gotp"
)

const (
	CLOTP = "clotp"
)

type Keys interface {
	ListOTPs() ([]KeyringKey, error)
	AddKey(name string, otpUri string) error
	GetByName(name string) (*KeyringItem, error)
	GetKeyByName(name string) (*KeyringKey, error)
	RemoveByName(name string) error
	UpdateKey(name string, data *KeyringItem) error
}

type KeyringKeys struct {
	Keys
	ring keyring.Keyring
}

func NewKeys() (*KeyringKeys, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName:              CLOTP,
		KeychainName:             CLOTP,
		KWalletAppID:             CLOTP,
		KWalletFolder:            CLOTP,
		WinCredPrefix:            CLOTP,
		LibSecretCollectionName:  CLOTP,
		KeychainTrustApplication: true,
	})
	if err != nil {
		return nil, err
	}
	return &KeyringKeys{
		ring: ring,
	}, nil
}

func dataToOtpKey(item keyring.Item) (*KeyringItem, error) {
	if len(item.Data) != 0 {
		params, err := parseOtpParams(string(item.Data))
		if err != nil {
			return nil, err
		}
		key, err := keyFromKeyringItem(&item)
		key.Name = item.Label
		if err != nil {
			return nil, err
		}
		return &KeyringItem{Key: key, OTP: params}, nil
	}
	return nil, nil
}

func (k *KeyringKeys) ListOTPs() ([]KeyringKey, error) {
	result := make([]KeyringKey, 0)
	keys, err := k.ring.Keys()
	if err != nil {
		return nil, err
	}
	for _, keySt := range keys {
		meta, err := k.ring.GetMetadata(keySt)
		if err != nil {
			return nil, err
		}
		key, err := keyFromKeyringItem(meta.Item)
		if err != nil {
			return nil, err
		}
		result = append(result, *key)
	}
	return result, err
}

func getType(otp gotp.OTP) OtpType {
	switch otp.(type) {
	case *gotp.HOTP:
		return HOTP
	case *gotp.TOTP:
		return TOTP
	default:
		return Unknown
	}
}

func (k *KeyringKeys) findByLabel(label string) (*keyring.Item, error) {
	keys, err := k.ring.Keys()
	if err != nil {
		return nil, err
	}
	for _, keyStr := range keys {
		meta, err := k.ring.GetMetadata(keyStr)
		if err != nil {
			return nil, err
		}
		if meta.Label == label {
			return meta.Item, nil
		}
	}
	return nil, nil
}

func (k *KeyringKeys) AddKey(name string, otpUri string) error {
	otp, err := gotp.OTPFromUri(otpUri)
	if err != nil {
		return err
	}
	otpParams, err := ParamsFromOTP(otp.OTP)
	if err != nil {
		return err
	}
	var label string
	if name == "" {
		label = otp.GetLabelRepr()
	} else {
		label = name
	}
	key := KeyringKey{
		Name:    label,
		Type:    getType(otp.OTP),
		Account: otp.Account,
		Issuer:  otp.Issuer,
	}
	return k.saveItem(&KeyringItem{
		Key: &key,
		OTP: otpParams,
	})
}

func (k *KeyringKeys) saveItem(item *KeyringItem) error {
	keyStr, err := item.Key.toKey()
	if err != nil {
		return err
	}
	paramsStr, err := item.OTP.asString()
	if err != nil {
		return err
	}
	kitem := keyring.Item{
		Key:                         keyStr,
		Data:                        []byte(paramsStr),
		Label:                       item.Key.Name,
		KeychainNotTrustApplication: true,
	}
	return k.ring.Set(kitem)
}

func (k *KeyringKeys) GetByName(name string) (*KeyringItem, error) {
	item, err := k.findByLabel(name)
	if err != nil {
		return nil, err
	}
	if item != nil {
		data, err := k.ring.Get(item.Key)
		if err != nil {
			return nil, err
		}
		return dataToOtpKey(data)
	} else {
		return nil, fmt.Errorf("OTP code %s not found", name)
	}
}

func (k *KeyringKeys) GetKeyByName(name string) (*KeyringKey, error) {
	item, err := k.findByLabel(name)
	if err != nil {
		return nil, err
	}
	if item != nil {
		data, err := k.ring.GetMetadata(item.Key)
		if err != nil {
			return nil, err
		}
		return keyFromKeyringItem(data.Item)
	} else {
		return nil, fmt.Errorf("OTP code %s not found", name)
	}
}

func (k *KeyringKeys) RemoveByName(name string) error {
	item, err := k.findByLabel(name)
	if err != nil {
		return err
	}
	if item == nil {
		return nil
	}
	return k.ring.Remove(item.Key)
}

func (k *KeyringKeys) UpdateKey(name string, key *KeyringItem) error {
	item, err := k.findByLabel(name)
	if err != nil {
		return err
	}
	if item == nil {
		return fmt.Errorf("OTP code %s not found", name)
	}
	err = k.RemoveByName(name)
	if err != nil {
		return err
	}
	return k.saveItem(key)
}
