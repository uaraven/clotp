package keyrings

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/99designs/keyring"
	"github.com/uaraven/gotp"
)

type OtpType int

const (
	Unknown OtpType = iota
	TOTP
	HOTP
)

// OtpParams contains all the parameters needed to calculate OTP
// This structure is marshalled to JSON, encoded to base64 and stored in keyring as a secret
type OtpParams struct {
	Hash      crypto.Hash `json:"hash"`
	Secret    []byte      `json:"secret"`
	Digits    int         `json:"digits"`
	StartTime int64       `json:"startTime,omitempty"`
	TimeStep  int         `json:"timeStep,omitempty"`
}

type KeyringKey struct {
	Name    string  `json:"-"`
	Type    OtpType `json:"type"`
	Account string  `json:"account"`
	Issuer  string  `json:"issuer"`
	Counter int64   `json:"counter,omitempty"`
}

type KeyringItem struct {
	Key *KeyringKey
	OTP *OtpParams
}

func parseOtpParams(otp string) (*OtpParams, error) {
	data, err := base64.StdEncoding.DecodeString(otp)
	if err != nil {
		return nil, err
	}
	var otpParam OtpParams
	err = json.Unmarshal(data, &otpParam)
	if err != nil {
		return nil, err
	}
	return &otpParam, nil
}

func ParamsFromUri(uri string) (*OtpParams, error) {
	otp, err := gotp.OTPFromUri(uri)
	if err != nil {
		return nil, err
	}
	return ParamsFromOTP(otp.OTP)
}

func ParamsFromOTP(otp gotp.OTP) (*OtpParams, error) {
	var out OtpParams
	out.Secret = otp.GetSecret()
	out.Hash = otp.GetHash()
	out.Digits = otp.GetDigits()
	switch otpt := otp.(type) {
	case *gotp.TOTP:
		// out.Type = TOTP
		out.StartTime = otpt.GetStartTime()
		out.TimeStep = otpt.GetTimeStep()
	case *gotp.HOTP:
		// out.Type = HOTP
	default:
		return nil, fmt.Errorf("unsupported OTP type: %T", otp)
	}
	return &out, nil
}

func (ki KeyringItem) AsOTP() (gotp.OTP, error) {
	otp := ki.OTP
	if ki.Key.Type == TOTP {
		return gotp.NewTOTPHash(otp.Secret, otp.Digits, otp.TimeStep, otp.StartTime, otp.Hash), nil
	} else if ki.Key.Type == HOTP {
		return gotp.NewHOTPHash(otp.Secret, ki.Key.Counter, otp.Digits, gotp.DefaultTransactionOffset, otp.Hash), nil
	} else {
		return nil, fmt.Errorf("unsupported OTP type: %d", ki.Key.Type)
	}
}

func (or *OtpParams) asString() (string, error) {
	data, err := json.Marshal(or)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func KeyFromOTPData(data gotp.OTPKeyData) (*KeyringKey, error) {
	var otype OtpType
	switch data.OTP.(type) {
	case *gotp.HOTP:
		otype = HOTP
	case *gotp.TOTP:
		otype = TOTP
	default:
		return nil, fmt.Errorf("unknown OTP type: %T", data.OTP)
	}
	return &KeyringKey{
		Type:    otype,
		Issuer:  data.Issuer,
		Account: data.Account,
	}, nil
}

func NewKey(otype OtpType, issuer string, account string) *KeyringKey {
	return &KeyringKey{
		Type:    otype,
		Issuer:  issuer,
		Account: account,
	}
}

func keyFromKeyringItem(keyItem *keyring.Item) (*KeyringKey, error) {
	if keyItem == nil {
		return nil, fmt.Errorf("invalid keyring key")
	}
	data, err := base64.StdEncoding.DecodeString(keyItem.Key)
	if err != nil {
		return nil, err
	}
	var ki KeyringKey
	err = json.Unmarshal(data, &ki)
	if err != nil {
		return nil, err
	}
	ki.Name = keyItem.Label
	return &ki, nil
}

func (kk *KeyringKey) toKey() (string, error) {
	data, err := json.Marshal(kk)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (kk KeyringKey) GetTypeString() string {
	switch kk.Type {
	case HOTP:
		return "HOTP"
	case TOTP:
		return "TOTP"
	default:
		return "Unknown"
	}
}

func (kk KeyringKey) GetLabel() string {
	if kk.Issuer != "" {
		return fmt.Sprintf("%s (%s)", kk.Account, kk.Issuer)
	} else {
		return kk.Account
	}
}
