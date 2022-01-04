package main

import (
	"crypto"
	"fmt"
	"net/url"

	"github.com/uaraven/clotp/migration"
	"github.com/uaraven/gotp"
)

const migrationScheme = "otpauth-migration"

func otpFromMigrationUri(uri string) ([]gotp.OTPKeyData, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if u.Scheme != migrationScheme {
		return nil, fmt.Errorf("unsupported URL scheme: %s, expected '%s'", u.Scheme, migrationScheme)
	}
	payload, err := migration.UnmarshalURL(uri)
	if err != nil {
		return nil, err
	}
	result := make([]gotp.OTPKeyData, 0)
	for _, params := range payload.OtpParameters {
		otp, err := otpFromParameters(params)
		if err != nil {
			return nil, err
		}
		okd := gotp.OTPKeyData{
			Issuer: params.Issuer,
			Label:  params.Name,
			OTP:    otp,
		}
		result = append(result, okd)
	}
	return result, nil
}

func otpFromParameters(params *migration.Payload_OtpParameters) (gotp.OTP, error) {
	switch params.Type {
	case migration.Payload_OTP_TYPE_TOTP:
		return totpFromParameters(params)
	case migration.Payload_OTP_TYPE_HOTP:
		return hotpFromParameters(params)
	default:
		return nil, fmt.Errorf("unsupported OTP algorithm")
	}
}

func totpFromParameters(params *migration.Payload_OtpParameters) (gotp.OTP, error) {
	hash, err := hashFromAlgorithm(params.Algorithm)
	if err != nil {
		return nil, err
	}
	return gotp.NewTOTPHash(
		params.Secret,
		decodeDigits(params.Digits),
		gotp.DefaultTimeStep,
		0,
		hash,
	), nil
}

func hotpFromParameters(params *migration.Payload_OtpParameters) (gotp.OTP, error) {
	hash, err := hashFromAlgorithm(params.Algorithm)
	if err != nil {
		return nil, err
	}
	return gotp.NewHOTPHash(
		params.Secret,
		int64(params.Counter),
		decodeDigits(params.Digits),
		0,
		hash,
	), nil
}

func hashFromAlgorithm(algo migration.Payload_Algorithm) (crypto.Hash, error) {
	switch algo {
	case migration.Payload_ALGORITHM_SHA1:
		return crypto.SHA1, nil
	case migration.Payload_ALGORITHM_SHA256:
		return crypto.SHA256, nil
	case migration.Payload_ALGORITHM_SHA512:
		return crypto.SHA3_512, nil
	case migration.Payload_ALGORITHM_UNSPECIFIED:
		return crypto.SHA1, nil
	default:
		return 0, fmt.Errorf("unsupported hash algorithm: %s", migration.Payload_Algorithm_name[int32(algo)])
	}
}

func decodeDigits(digits migration.Payload_DigitCount) int {
	switch digits {
	case migration.Payload_DIGIT_COUNT_EIGHT:
		return 8
	default:
		return 6
	}
}
