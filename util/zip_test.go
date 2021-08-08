package util

import (
	"testing"
)

func Test_zippie_ZipWriter(t *testing.T) {
	z := NewZippie("store")
	writer, _, err := z.Zip()
	if err != nil {
		t.Error(err)
	}
	t.Log(writer)

}
