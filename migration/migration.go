package migration

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"google.golang.org/protobuf/proto"
)

// ErrUnknown scheme or host
var ErrUnknown = errors.New("unknown")

// Data extracts data part from URL string
func Data(link string) ([]byte, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "otpauth-migration" {
		return nil, fmt.Errorf("scheme %s: %w", u.Scheme, ErrUnknown)
	}
	if u.Host != "offline" {
		return nil, fmt.Errorf("host %s: %w", u.Host, ErrUnknown)
	}
	data := u.Query().Get("data")
	// fix spaces back to plus sign
	data = strings.ReplaceAll(data, " ", "+")
	return base64.StdEncoding.DecodeString(data)
}

// Unmarshal otpauth-migration data
func Unmarshal(data []byte) (*Payload, error) {
	var p Payload
	if err := proto.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// UnmarshalURL decodes otpauth-migration from URL
func UnmarshalURL(link string) (*Payload, error) {
	data, err := Data(link)
	if err != nil {
		return nil, err
	}
	return Unmarshal(data)
}

func MarshalToURL(p *Payload) (string, error) {
	data, err := proto.Marshal(p)
	if err != nil {
		return "", err
	}
	base64data := base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(data)

	url := url.URL{
		Scheme:   "otpauth-migration",
		Host:     "offline",
		RawQuery: "data=" + url.QueryEscape(base64data),
	}
	return url.String(), nil
}
