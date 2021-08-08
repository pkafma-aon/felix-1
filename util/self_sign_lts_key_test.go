package util

import "testing"

func TestGenerateLtsCertKeyPinPair(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"one", args{"self.sign"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenerateLtsCertKeyPinPair(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("GenerateLtsCertKeyPinPair() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
