package utils

import (
	"encoding/base64"
)

// Base64Encrypt Base64加密，传入对应的text文本
func Base64Encrypt(text string) (encryptionText string) {
	strBytes := []byte(text)
	encryptionText = base64.StdEncoding.EncodeToString(strBytes)
	return
}

// Base64Decrypt Base64解密，传入base64加密后的文本
func Base64Decrypt(base64Text string) (decryptionText string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(base64Text)
	//if err != nil {
	//	common.LogErrorf("Text Decrypt Failed", logrus.Fields{"err": err, "base64Text": base64Text})
	//}
	decryptionText = string(decoded)
	return
}
