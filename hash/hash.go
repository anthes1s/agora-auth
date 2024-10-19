package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

type hashConfig struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var config = hashConfig{
	memory:      64 * 1024,
	iterations:  2,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

// NOTE: Stolen from: https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
func Generate(password string) (string, error) {
	salt, err := generateSalt(config.saltLength)
	if err != nil {
		log.Printf("Error: Failed to generate salt: %v", err)
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		config.iterations,
		config.memory,
		config.parallelism,
		config.keyLength,
	)

	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		config.memory,
		config.iterations,
		config.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func Verify(password string, encoded string) (bool, error) {
	config, salt, hash, err := decodeHash(encoded)
	if err != nil {
		log.Printf("Error: Failed to decode hash - %v\n", err)
		return false, nil
	}

	clientHash := argon2.IDKey(
		[]byte(password),
		salt,
		config.iterations,
		config.memory,
		config.parallelism,
		config.keyLength,
	)

	hashLength := int32(len(hash))
	givenHashLength := int32(len(clientHash))

	if subtle.ConstantTimeEq(hashLength, givenHashLength) == 0 {
		return false, nil
	}
	if subtle.ConstantTimeCompare(hash, clientHash) == 1 {
		return true, nil
	}

	return false, nil
}

func generateSalt(amount uint32) ([]byte, error) {
	b := make([]byte, amount)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error: Failed to generate salt: %v", err)
		return nil, err
	}

	return b, nil
}

func decodeHash(encoded string) (*hashConfig, []byte, []byte, error) {
	vals := strings.Split(encoded, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("Invalid hash")
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, errors.New("Invalid argon variant")
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("Invalid argon version")
	}

	params := &hashConfig{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.saltLength = uint32(len(salt))

	key, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.keyLength = uint32(len(key))

	return params, salt, key, nil
}
