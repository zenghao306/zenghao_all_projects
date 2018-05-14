package common

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"
)

type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	Sign(src []byte, hash crypto.Hash) ([]byte, error)
	Verify(src []byte, sign []byte, hash crypto.Hash) error
}

type pkcsClient struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (this *pkcsClient) Encrypt(plaintext []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, this.publicKey, plaintext)
}
func (this *pkcsClient) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, this.privateKey, ciphertext)
}

func (this *pkcsClient) Sign(src []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, this.privateKey, hash, hashed)
}

func (this *pkcsClient) Verify(src []byte, sign []byte, hash crypto.Hash) error {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(this.publicKey, hash, hashed, sign)
}

//默认客户端，pkcs8私钥格式，pem编码
func NewDefault(publicKey string) (Cipher, error) {
	blockPub, _ := pem.Decode([]byte(publicKey))
	if blockPub == nil {
		return nil, errors.New("public key error")
	}

	return newPublicKey(blockPub.Bytes)
}

func newPublicKey(publicKey []byte) (Cipher, error) {

	pubKey, err := GenPubKey(publicKey)
	if err != nil {
		return nil, err
	}
	return &pkcsClient{privateKey: nil, publicKey: pubKey}, nil
}

func GenPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

var cipher Cipher

func Init() {
	fmt.Printf("\ninit走到了\n")
	client, err := NewDefault(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwov2pIZv7apU9OaGG3NfyDTSVek622bIH+MvepDhrVlC8GUM55aulsw817KFyfH3Wc1jxcVShxErQmUBiUiOxs8VEMdJd1CYvsmLDu8mkODBChq2yanoSUJ0x03OuDjClLT6hp2MJZ8rq0ae18K9pLQvzPWmJSwToDRtbZw62YkwYHub61Pu/IHpGf/MkooVo+zPQenSt/DKUj9xxOmQa16Ify5u8rlqAe3JqMuaChc2aM7vQZvCwaHKSVor5chc6T1MBNKzG9AnQ4MiH/jwkn/5K/dCUvQ1C9IdXre9r37YPxc0RJYAtYd/PCjH7+L0Np7Z+d83htANudcu9tkEmwIDAQAB
-----END PUBLIC KEY-----`)

	if err != nil {
		fmt.Println(err)
	}

	cipher = client
}

func Test_DefaultClient(t *testing.T) {

	cp, err := cipher.Encrypt([]byte("测试加密解密"))
	if err != nil {
		t.Error(err)
	}
	cpStr := base64.URLEncoding.EncodeToString(cp)

	fmt.Println(cpStr)

	ppBy, err := base64.URLEncoding.DecodeString(cpStr)
	if err != nil {
		t.Error(err)
	}
	pp, err := cipher.Decrypt(ppBy)

	fmt.Println(string(pp))
}

func TestVerify(strBeforSig, sigStr string) bool {
	Init()
	//src := "测试签名验签"
	//src := "app_id=2015102700040153&body=大乐透2.1&buyer_id=2088102116773037&charset=utf-8&gmt_close=2016-07-19 14:10:46&gmt_create=2016-07-19 14:10:44&gmt_payment=2016-07-19 14:10:47&notify_id=4a91b7a78a503640467525113fb7d8bg8e&notify_time=2016-07-19 14:10:49&notify_type=trade_status_sync&out_trade_no=0719141034-6418&refund_fee=0.00&seller_id=2088102119685838&subject=大乐透2.1&total_amount=2.00&trade_no=2016071921001003030200089909&trade_status=TRADE_SUCCESS&version=1.0"
	//signBytes, err := cipher.Sign([]byte(src), crypto.SHA256)
	//fmt.Printf("\n126")
	//if err != nil {
	//	fmt.Printf("\n出错了1111")
	//	return
	//}

	//sign := hex.EncodeToString(signBytes)
	//fmt.Println(sign)
	//fmt.Printf("\n128")

	signB, _ := hex.DecodeString(sigStr)

	fmt.Printf("\nstrBeforSig:%s\n", strBeforSig)
	fmt.Printf("\nsigStr:%s\n", sigStr)

	errV := cipher.Verify([]byte(strBeforSig), signB, crypto.SHA256)
	if errV != nil {
		fmt.Printf("\nTestVerify() 2222")
		return false
	}
	fmt.Println("verify success")
	return true
}
