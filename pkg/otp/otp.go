package otp

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	optUpdateInterval = 600
)

func getHotp(key []byte, seed []byte) (string, error) {
	hash := hmac.New(sha1.New, key)
	if _, err := hash.Write(seed); err != nil {
		return "", err
	}

	h := hash.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	o := (h[19] & 15)

	var header uint32
	//Get 32 bit chunk from hash starting at the o
	r := bytes.NewReader(h[o : o+4])
	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		return "", err
	}

	//Ignore most significant bits as per RFC 4226.
	//Takes division from one million to generate a remainder less than < 7 digits
	h12 := (int(header) & 0x7fffffff) % 1000000

	return fmt.Sprintf("%06d", h12), nil
}

func GetTotp(key []byte) (prev string, curr string, err error) {
	now := uint64(time.Now().Unix() / optUpdateInterval)

	seed := make([]byte, 8)
	binary.BigEndian.PutUint64(seed, now-1)
	if prev, err = getHotp(key, seed); err != nil {
		return "", "", err
	}
	binary.BigEndian.PutUint64(seed, now)
	if curr, err = getHotp(key, seed); err != nil {
		return prev, "", err
	}

	return prev, curr, err
}
