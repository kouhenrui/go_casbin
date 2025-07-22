package pkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// RSAEncryptor 接口
type RSAEncryptor interface {
	MakeSalt() string                                        // 16位密码加密盐
	Rand6String() string                                     // 随机6位密钥
	RandAllString() string                                   // 随机16位密钥
	Bcrypt(pwd string) (string, error)                       // 加密
	BcryptCheck(pwd string, hash string) (bool, error)       // 解密
	EnPwdCode(pwd string, pwdKey string) (string, error)     // 加密
	DePwdCode(pwd string, pwdKey string) (string, error)     // 解密
	EncryptWithPublicKey(plaintext []byte) (string, error)   // 公钥加密
	DecryptWithPrivateKey(ciphertext string) ([]byte, error) // 私钥解密
	EncryptAES(plaintext []byte, key []byte) ([]byte, []byte, error) // aes256加密
	DecryptAES(ciphertext []byte, key []byte, iv []byte) (string, error) // aes256解密
}

type DefaultEncryptor struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func (e *DefaultEncryptor) MakeSalt() string {
	return randString(16)
}

func (e *DefaultEncryptor) Rand6String() string {
	return randString(6)
}

func (e *DefaultEncryptor) RandAllString() string {
	return randString(16)
}

func (e *DefaultEncryptor) EnPwdCode(pwd string, pwdKey string) (string, error) {
	key := []byte(pwdKey)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("key must be 16, 24, or 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintext := []byte(pwd)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext = append(plaintext, padtext...)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(append(iv, ciphertext...)), nil
}

func (e *DefaultEncryptor) DePwdCode(pwd string, pwdKey string) (string, error) {
	key := []byte(pwdKey)
	data, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	padding := int(ciphertext[len(ciphertext)-1])
	return string(ciphertext[:len(ciphertext)-padding]), nil
}

func (e *DefaultEncryptor) EncryptWithPublicKey(plaintext []byte) (string, error) {
	if e.PublicKey == nil {
		return "", errors.New("public key not set")
	}
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, e.PublicKey, plaintext)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *DefaultEncryptor) DecryptWithPrivateKey(ciphertext string) ([]byte, error) {
	if e.PrivateKey == nil {
		return nil, errors.New("private key not set")
	}
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, e.PrivateKey, data)
}

func (e *DefaultEncryptor) EncryptAES(plaintext []byte, key []byte) ([]byte, []byte, error) {
	if len(key) != 32 {
		return nil, nil, errors.New("AES-256 key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext = append(plaintext, padtext...)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, iv, nil
}

func (e *DefaultEncryptor) DecryptAES(ciphertext []byte, key []byte, iv []byte) (string, error) {
	if len(key) != 32 {
		return "", errors.New("AES-256 key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	padding := int(ciphertext[len(ciphertext)-1])
	return string(ciphertext[:len(ciphertext)-padding]), nil
}

func (e *DefaultEncryptor) Bcrypt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (e *DefaultEncryptor) BcryptCheck(pwd string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	if err != nil {
		return false, err
	}
	return true, nil
}

// 工具函数
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[num.Int64()]
	}
	return string(result)
}

// PEM 解析工具
func ParseRSAPublicKeyFromPEM(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func ParseRSAPrivateKeyFromPEM(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
} 