package cli

import (
	"fmt"
	"strings"

	"github.com/uaraven/clotp/keyrings"
)

type ListCmd struct {
}

func List(cmd *ListCmd, keys keyrings.Keys) (string, error) {
	otpKeys, err := keys.ListOTPs()
	if err != nil {
		return "", err
	}
	if len(otpKeys) == 0 {
		return "No stored OTPs", nil
	} else {
		var result strings.Builder
		for _, otpKey := range otpKeys {
			result.WriteString(fmt.Sprintf("%s:  %s\n", otpKey.GetTypeString(), otpKey.GetLabel()))
		}
		return result.String(), nil
	}
}
