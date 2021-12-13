package methods

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"log"
)

// Sign is method used as HMAC-SHA1 signer
func Sign(signingKey, message string) string {
	// generate new hash mac
	mac := hmac.New(sha1.New, []byte(signingKey))
	number, err := mac.Write([]byte(message))
	if err != nil {
		log.Println(err)
		return ""
	}
	log.Println(number)
	signatureBytes := mac.Sum(nil)
	// return encoded string using base64 encoding
	return base64.StdEncoding.EncodeToString(signatureBytes)
}
