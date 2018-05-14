package common

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

type Type int64

const (
	PKCS1 Type = iota
	PKCS8
)

func Rsa256Signal(sigStr, privateKey string) string {
	hash := sha256.New()
	hash.Write([]byte(sigStr))
	rsaString := hash.Sum(nil)

	blockPri, _ := pem.Decode([]byte(privateKey))
	if blockPri == nil {
		fmt.Printf("pem.Decode err.\r\n")
		return ""
	}

	priv, err := genPriKey(blockPri.Bytes, PKCS8)
	if err != nil {
		fmt.Printf("genPriKey err %v.\r\n", err)
		return ""
	}

	signCode, _ := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, rsaString)

	return base64.URLEncoding.EncodeToString(signCode)
}

func genPriKey(privateKey []byte, privateKeyType Type) (*rsa.PrivateKey, error) {
	var priKey *rsa.PrivateKey
	var err error
	switch privateKeyType {
	case PKCS1:
		{
			priKey, err = x509.ParsePKCS1PrivateKey(privateKey)
			if err != nil {
				return nil, err
			}
		}
	case PKCS8:
		{
			prkI, err := x509.ParsePKCS8PrivateKey(privateKey)
			if err != nil {
				return nil, err
			}
			priKey = prkI.(*rsa.PrivateKey)
		}
	default:
		{
			return nil, errors.New("unsupport private key type")
		}
	}

	return priKey, nil
}

func Rsa256PublicKeySignal(sigStr, publicKey string) string {
	hash := sha256.New()
	hash.Write([]byte(sigStr))
	rsaString := hash.Sum(nil)

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return ""
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return ""
	}
	pub := pubInterface.(*rsa.PublicKey)
	signCode, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, rsaString)

	return base64.URLEncoding.EncodeToString(signCode)
}

func genPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func NewPublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	pubKey, err := genPubKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

///////////////////////////////////////////////
