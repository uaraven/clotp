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
		var cols [4]int = [4]int{8, 4, 7, 6}
		for _, otpKey := range otpKeys {
			if len(otpKey.Name) > cols[1] {
				cols[1] = len(otpKey.Name)
			}
			if len(otpKey.Account) > cols[2] {
				cols[2] = len(otpKey.Account)
			}
			if len(otpKey.Issuer) > cols[3] {
				cols[3] = len(otpKey.Issuer)
			}
		}

		fmt.Printf("%*s %*s %*s %*s\n", cols[0], "OTP Type", cols[1], "Name", cols[2], "Account", cols[3], "Issuer")
		fmt.Printf("%s %s %s %s\n",
			strings.Repeat("-", cols[0]),
			strings.Repeat("-", cols[1]),
			strings.Repeat("-", cols[2]),
			strings.Repeat("-", cols[3]),
		)
		var result strings.Builder
		for _, otpKey := range otpKeys {
			result.WriteString(fmt.Sprintf("%*s %*s %*s %*s\n",
				cols[0],
				otpKey.GetTypeString(),
				cols[1],
				otpKey.Name,
				cols[2],
				otpKey.Account,
				cols[3],
				otpKey.Issuer))
		}
		return result.String(), nil
	}
}
