package crypto

import (
	"crypto/rand"

	argon "golang.org/x/crypto/argon2"
)

const (
	argonSaltLen = 16
	argonTime    = 1
	argonMemory  = 1 << 21 // 2 Gib
	argonThreads = 4
	argonKeyLen  = 32
)

func Argon2(password []byte, salt []byte) []byte {
	return argon.IDKey(password, salt, argonTime, argonMemory, argonThreads, argonKeyLen)
}

func Argon2Hash(password []byte) ([]byte, []byte, error) {
	salt := make([]byte, argonSaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, nil, err
	}

	key := Argon2(password, salt)
	return key, salt, nil
}
