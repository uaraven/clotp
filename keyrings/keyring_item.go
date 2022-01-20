package keyrings

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"

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
	Type      OtpType     `json:"type"`
	Hash      crypto.Hash `json:"hash"`
	Secret    []byte      `json:"secret"`
	Digits    int         `json:"digits"`
	Counter   int64       `json:"counter,omitempty"`
	StartTime int64       `json:"startTime,omitempty"`
	TimeStep  int         `json:"timeStep,omitempty"`
}

type KeyringKey struct {
	Type    OtpType `json:"type"`
	Account string  `json:"account"`
	Issuer  string  `json:"issuer"`
}

type KeyringItem struct {
	Key *KeyringKey
	OTP *OtpParams
}

func ParseOtpParams(otp string) (*OtpParams, error) {
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
		out.Type = TOTP
		out.StartTime = otpt.GetStartTime()
		out.TimeStep = otpt.GetTimeStep()
	case *gotp.HOTP:
		out.Type = HOTP
		out.Counter = otpt.GetCounter()
	default:
		return nil, fmt.Errorf("unsupported OTP type: %T", otp)
	}
	return &out, nil
}

func (or *OtpParams) AsOTP() (gotp.OTP, error) {
	if or.Type == TOTP {
		return gotp.NewTOTPHash(or.Secret, or.Digits, or.TimeStep, or.StartTime, or.Hash), nil
	} else if or.Type == HOTP {
		return gotp.NewHOTPHash(or.Secret, or.Counter, or.Digits, -1, or.Hash), nil
	} else {
		return nil, fmt.Errorf("unsupported OTP type: %d", or.Type)
	}
}

func (or *OtpParams) AsString() (string, error) {
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
		Account: data.Label,
	}, nil
}

func NewKey(otype OtpType, issuer string, label string) *KeyringKey {
	return &KeyringKey{
		Type:    otype,
		Issuer:  issuer,
		Account: label,
	}
}

func KeyFromString(keyS string) (*KeyringKey, error) {
	data, err := base64.StdEncoding.DecodeString(keyS)
	if err != nil {
		return nil, err
	}
	var ki KeyringKey
	err = json.Unmarshal(data, &ki)
	if err != nil {
		return nil, err
	}
	return &ki, nil
}

func (kk *KeyringKey) ToKey() (string, error) {
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
