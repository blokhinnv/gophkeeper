package encrypt

import "encoding/base64"

func EncryptString(text, key string) (string, error) {
	res, err := EncryptBytes([]byte(text), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func DecryptString(encryptedText, key string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	res, err := DecryptBytes(b, key)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
