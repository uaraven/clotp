package cli

import (
	"fmt"

	"github.com/uaraven/clotp/keyrings"
)

type RemoveCmd struct {
	Name string `arg:"positional,required" help:"Name of the code to remove"`
}

func Remove(cmd *RemoveCmd, keys keyrings.Keys) (string, error) {
	var output string
	if cmd.Name == "" {
		return "", fmt.Errorf("OTP code name must be specified")
	}
	output = fmt.Sprintf("Removing code with name '%s'\n", cmd.Name)
	return output, keys.RemoveByName(cmd.Name)
}
