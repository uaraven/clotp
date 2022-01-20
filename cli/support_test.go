package cli

import (
	"fmt"

	"github.com/uaraven/clotp/keyrings"
	"github.com/uaraven/gotp"
)

type mockKeys struct {
	keyrings.Keys
	keys map[keyrings.KeyringKey]keyrings.OtpParams
}

func (mk *mockKeys) ListOTPs() ([]keyrings.KeyringKey, error) {
	result := make([]keyrings.KeyringKey, 0)
	for k := range mk.keys {
		result = append(result, k)
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
	mk.keys[*key] = *data
	return nil
}

func (mk *mockKeys) GetByName(name string) (*keyrings.KeyringItem, error) {
	for key, item := range mk.keys {
		if key.Account == name {
			return &keyrings.KeyringItem{
				Key: &key,
				OTP: &item,
			}, nil
		}
	}
	return nil, fmt.Errorf("OTP with name %s not found", name)
}

func (mk *mockKeys) RemoveByName(name string) error {
	for key := range mk.keys {
		if key.Account == name {
			delete(mk.keys, key)
			return nil
		}
	}
	return fmt.Errorf("OTP with name %s not found", name)
}

func NewMockKeys() *mockKeys {
	return &mockKeys{
		keys: make(map[keyrings.KeyringKey]keyrings.OtpParams),
	}
}
