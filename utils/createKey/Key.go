package createKey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

type De struct {
	Str string
	Err error
}

// GenerateRSAKeys 生成 RSA 密钥并保存到文件
func GenerateRSAKeys() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 保存私钥
	privFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer privFile.Close()
	pem.Encode(privFile, &pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	// 保存公钥
	pubFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer pubFile.Close()
	pubASN1 := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubASN1})

	return nil
}

// LoadPrivateKey 从文件中加载私钥
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// DecryptPassword 解密密码
func DecryptPassword(privateKey *rsa.PrivateKey, encryptedPassword string) De {

	cipherText, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return De{
			Str: "",
			Err: err,
		}
	}

	decryptedPassword, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		return De{
			Str: "",
			Err: err,
		}
	}

	return De{Str: string(decryptedPassword), Err: err}
}
