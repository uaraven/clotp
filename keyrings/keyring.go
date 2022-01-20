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
	RemoveByName(name string) error
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
		params, err := ParseOtpParams(string(item.Data))
		if err != nil {
			return nil, err
		}
		key, err := KeyFromString(item.Key)
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
		key, err := KeyFromString(keySt)
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
		label = otp.Label
	} else {
		label = name
	}
	key := KeyringKey{
		Type:    getType(otp.OTP),
		Account: label,
		Issuer:  otp.Issuer,
	}
	keyStr, err := key.ToKey()
	if err != nil {
		return err
	}
	paramsStr, err := otpParams.AsString()
	if err != nil {
		return err
	}
	kitem := keyring.Item{
		Key:   keyStr,
		Data:  []byte(paramsStr),
		Label: label,
	}
	return k.ring.Set(kitem)
}

func (k *KeyringKeys) GetByName(name string) (*KeyringItem, error) {
	keys, err := k.ring.Keys()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		kk, err := KeyFromString(key)
		if err != nil {
			return nil, err
		}
		if kk.Account == name {
			item, err := k.ring.Get(key)
			if err != nil {
				return nil, err
			}
			return dataToOtpKey(item)
		}
	}
	return nil, fmt.Errorf("cannot find code with name=%s", name)
}

func (k *KeyringKeys) RemoveByName(name string) error {
	keys, err := k.ring.Keys()
	if err != nil {
		return err
	}
	for _, key := range keys {
		kk, err := KeyFromString(key)
		if err != nil {
			return err
		}
		if kk.Account == name {
			// remove always fails, so we just clear the data for now
			// return k.ring.Remove(item.Key)
			item, err := k.ring.Get(key)
			if err != nil {
				return err
			}
			item.Data = make([]byte, 0)
			return k.ring.Set(item)
		}
	}
	return fmt.Errorf("cannot find key with name=%s", name)
}
