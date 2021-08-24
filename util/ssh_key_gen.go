package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
)

//SshKeyPairGenerate
//https://gist.github.com/devinodaniel/8f9b8a4f31573f428f29ec0e884e6673#file-gistfile1-txt-L29
func SshKeyPairGenerate() (publicKeyString, privateKeyString string, err error) {
	bitSize := 1024

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		return "", "", err
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", nil
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	return string(publicKeyBytes), string(privateKeyBytes), nil
}

// generatePrivateKey creates a RSA Private Name of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {

	// Private Name generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Name
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Name from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privateKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	return pubKeyBytes, nil
}
