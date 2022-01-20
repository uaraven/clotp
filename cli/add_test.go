package cli

import "testing"

func TestAddTOTP(t *testing.T) {
	cmd := AddCmd{
		Uri:  "otpauth://totp/vendor:label?secret=NNSXS&issuer=vendor",
		Name: "",
	}
	keys := NewMockKeys()
	_, err := Add(&cmd, keys)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(keys.keys) != 1 {
		t.Errorf("expected 1 key in the keyring")
	}
}

func TestAddHOTP(t *testing.T) {
	cmd := AddCmd{
		Uri:  "otpauth://hotp/vendor:label?secret=NNSXS&issuer=vendor&counter=1",
		Name: "",
	}
	keys := NewMockKeys()
	_, err := Add(&cmd, keys)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(keys.keys) != 1 {
		t.Errorf("expected 1 key in the keyring")
	}
}

func TestAddMigration(t *testing.T) {
	cmd := AddCmd{
		Uri:  "otpauth-migration://offline?data=ChwKBEtleTESBVRlc3QxGgdJc3N1ZXIxIAEoATACChwKBEtleTISBVRlc3QyGgdJc3N1ZXIyIAEoATAC",
		Name: "",
	}
	keys := NewMockKeys()
	_, err := Add(&cmd, keys)
	if err != nil {
		t.Error(err)
	}
	if len(keys.keys) != 2 {
		t.Errorf("expected 2 keys in the keyring")
	}
}
