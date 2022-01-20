package cli

import (
	"fmt"

	"github.com/uaraven/clotp/keyrings"
	"github.com/uaraven/gotp"
)

type mockKeys struct {
	keyrings.Keys
	keys map[string]keyrings.KeyringItem
}

func (mk *mockKeys) ListOTPs() ([]keyrings.KeyringKey, error) {
	result := make([]keyrings.KeyringKey, 0)
	for _, v := range mk.keys {
		result = append(result, *v.Key)
	}
	return result, nil
}

func (mk *mockKeys) AddKey(name string, otpUri string) error {
	otp, err := gotp.OTPFromUri(otpUri)
	if err != nil {
		return err
	}
	key, err := keyrings.KeyFromOTPData(*otp)
	if err != nil {
		return err
	}
	data, err := keyrings.ParamsFromOTP(otp.OTP)
	if err != nil {
		return err
	}
	mk.keys[name] = keyrings.KeyringItem{Key: key, OTP: data}
	return nil
}

func (mk *mockKeys) GetByName(name string) (*keyrings.KeyringItem, error) {
	item, ok := mk.keys[name]
	if !ok {
		return nil, fmt.Errorf("OTP with name %s not found", name)
	}
	return &item, nil
}

func (mk *mockKeys) RemoveByName(name string) error {
	_, ok := mk.keys[name]
	if !ok {
		return fmt.Errorf("OTP with name %s not found", name)
	} else {
		delete(mk.keys, name)
		return nil
	}
}

func (mk *mockKeys) UpdateKey(name string, data *keyrings.KeyringItem) error {
	mk.keys[name] = *data
	return nil
}

func NewMockKeys() *mockKeys {
	return &mockKeys{
		keys: make(map[string]keyrings.KeyringItem),
	}
}
