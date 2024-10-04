package utils

import (
	"encoding/base64"
)

func EncodeBase64(query string) string {
	return base64.StdEncoding.EncodeToString([]byte(query))
}

func DecodeBase64(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
