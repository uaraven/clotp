package cli

import (
	"github.com/uaraven/clotp/keyrings"
)

type CopyCmd struct {
	Name    string `arg:"positional,required" help:"Name of the OTP code"`
	Counter int64  `arg:"--counter" help:"Override counter for the HOTP code"`
}

func Copy(cmd *CopyCmd, keys keyrings.Keys) (string, error) {
	code := CodeCmd{
		Name:    cmd.Name,
		Counter: cmd.Counter,
		Copy:    true,
	}
	return Code(&code, keys)
}
