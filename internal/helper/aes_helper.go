package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	_ "fmt"
)

/*
*
你用 base64格式，使用postman 的时候，记得把 + / 从 AES-online 换成 - _
gw55ZcBQOW+lgUmjzCRyzA==
gw55ZcBQOW-lgUmjzCRyzA==
  ┌────────────────────┬───────────────────┬────────────────────────────────────┐
  │        编码        │      字符集       │              对应 PHP              │
  ├────────────────────┼───────────────────┼────────────────────────────────────┤
  │ base64.StdEncoding │ A-Z a-z 0-9 + / = │ base64_encode()                    │
  ├────────────────────┼───────────────────┼────────────────────────────────────┤
  │ base64.URLEncoding │ A-Z a-z 0-9 - _ = │ strtr(base64_encode(), '+/', '-_') │
  └────────────────────┴───────────────────┴────────────────────────────────────┘
*/

/*
Aes 128 需要 128 個 bit
然後 一個 英文字 可以 用一個 8bits（1byte） 的 utf8 表示
16 個英文字 剛好 等於 128bits
*/
type AesHelper struct {
	*AbstractHelper
}

func NewAesHelper(oAbstractHelper *AbstractHelper) *AesHelper {
	return &AesHelper{
		AbstractHelper: oAbstractHelper,
	}
}

func (oSelf *AesHelper) Encrypt(sText string, sKey string, sIv string) string {

	byText := []byte(sText)
	byKey := []byte(sKey)
	byIv := []byte(sIv)

	_byText := oSelf.pKCS7Padding(byText)
	ciphertext := make([]byte, len(_byText))
	mBlock, oErr := aes.NewCipher(byKey)
	if oErr != nil {
		panic(oErr)
	}

	oCipher := cipher.NewCBCEncrypter(mBlock, byIv)
	oCipher.CryptBlocks(ciphertext, _byText)

	sResult := base64.StdEncoding.EncodeToString(ciphertext)

	return sResult
}

func (oSelf *AesHelper) Decrypt(sText string, sKey string, sIv string) string {

	// 使用 base64.StdEncoding.DecodeString ->
	// 策略, 需要让前端传入 url版本的base64 (base64 原来版本的，用到了 http 协议的特殊符号 + /， 使用了会造成 http 解析异常)，
	// 然后 后端 改成 一般版本的base64
	// 使用标准 base64 解密

	// 使用 base64.UrlEncoding.DecodeString ->
	// 策略, 需要让前端传入 url版本的base64 (base64 原来版本的，用到了 http 协议的特殊符号 + / ， 使用了会造成 http 解析异常)，
	// 然后 后端 使用 url base64 版本 直接解密

	if sText == "" {
		return ""
	}

	oByteText, _ := base64.URLEncoding.DecodeString(sText)

	oByteKey := []byte(sKey)
	oByteIv := []byte(sIv)

	block, oErr := aes.NewCipher(oByteKey)
	if oErr != nil {
		panic(oErr)
	}

	ciphertext := make([]byte, len(oByteText))

	oCipher := cipher.NewCBCDecrypter(block, oByteIv)
	oCipher.CryptBlocks(ciphertext, oByteText)

	ciphertext = oSelf.pKCS7UnPadding(ciphertext)

	sResult := string(ciphertext)

	return sResult
}

func (oSelf *AesHelper) pKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (oSelf *AesHelper) pKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
