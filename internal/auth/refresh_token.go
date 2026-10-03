package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func MakeRefreshToken() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Print("error in filling the slice")
	}
	encodeStr := hex.EncodeToString(key)
	return encodeStr
}
