package util

import "testing"

func TestSshKeyPairGenerate(t *testing.T) {
	generate, privateKeyString, err := SshKeyPairGenerate()
	if err != nil {
		t.Error(err)
	}
	t.Log(generate)
	t.Log(privateKeyString)
}
