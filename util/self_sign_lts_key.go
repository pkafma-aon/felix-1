package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/big"
	"os"
	"time"
)

//publicPemBytes 暂存私有证书的公钥
//保存公钥到内存
var publicPemBytes []byte

//生成证书密钥对
//返回*tls.Certificate 做gRPC服务器启动时的参数
func GenerateTlsCert() (*tls.Certificate, error) {
	//1.- Generate private key:
	//随机种子
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	// Generate a pem block with the private key
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	//2.- Generate the certificate:
	tml := x509.Certificate{
		// you can add any attr that you need
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(20, 0, 0),
		// you have to generate a different serial number each execution
		SerialNumber: big.NewInt(9527),
		Subject: pkix.Name{
			CommonName:   "arrian.dev.mojotv.net", //可以自定义
			Organization: []string{"Cloud Security Management Platform", "云安全管理平台"},
		},
		BasicConstraintsValid: true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, &tml, &tml, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	// Generate a pem block with the certificate
	certPemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}
	publicPemBytes = pem.EncodeToMemory(certPemBlock)
	//2.1 write public.pem file
	tlsCert, err := tls.X509KeyPair(publicPemBytes, keyPem)
	return &tlsCert, err
}

//getLtsPublicKeyBytes 获取类型中的公钥 公钥 bytes
func getLtsPublicKeyBytes() string {
	if len(publicPemBytes) == 0 {
		logrus.Fatal("minion's SSL/TLS certificate has not generated, please make sure being called GenerateTlsCert first")
	}
	return string(publicPemBytes)
}

// KeyPairWithPin 返回 PEM证书 and PEM-Name 和SKPI(PIN码)
// 公共证书的指纹
func KeyPairWithPin() ([]byte, []byte, []byte, error) {
	bits := 4096
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("rsa.GenerateKey: %s", err)
	}

	tpl := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Issuer:                pkix.Name{Country: []string{"中国", "湖北省"}, Organization: []string{"EricZhou@mojotv.cn", "ByteGang"}},
		Subject:               pkix.Name{CommonName: "EricZhou@mojotv.cn-ByteGang", Locality: []string{"上海", "ByteGang"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(2, 0, 0),
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	derCert, err := x509.CreateCertificate(rand.Reader, &tpl, &tpl, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("x509.CreateCertificate: %s", err)
	}

	buf := &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derCert,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("pem.Encode: %s", err)
	}

	pemCert := buf.Bytes()

	buf = &bytes.Buffer{}
	err = pem.Encode(buf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("pem.Encode: %s", err)
	}
	pemKey := buf.Bytes()
	cert, err := x509.ParseCertificate(derCert)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("x509.ParseCertificate: %s", err)
	}

	pubDER, err := x509.MarshalPKIXPublicKey(cert.PublicKey.(*rsa.PublicKey))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("x509.MarshalPKIXPublicKey: %s", err)
	}
	sum := sha256.Sum256(pubDER)
	pin := make([]byte, base64.StdEncoding.EncodedLen(len(sum)))
	base64.StdEncoding.Encode(pin, sum[:])

	return pemCert, pemKey, pin, nil
}

func GenerateLtsCertKeyPinPair(name string) (err error) {
	certData, keyData, pinData, err := KeyPairWithPin()
	if err != nil {
		return err
	}
	m := map[string][]byte{
		"crt": certData,
		"key": keyData,
		"pin": pinData,
	}
	for k, data := range m {
		fn := fmt.Sprintf("https.%s.%s", name, k)
		f, err := os.Create(fn)
		if err != nil {
			return fmt.Errorf("创建文件失败:%s,   %s", fn, err)
		}
		defer f.Close()
		f.Write(data)
	}
	return nil
}
