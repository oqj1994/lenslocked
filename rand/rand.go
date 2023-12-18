package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("bytes:%w", err)
	}
	if nRead != n {
		panic(fmt.Errorf("didn't read enough bytes"))
	}
	return b, nil
}

func String(n int) (string, error) {
	bs, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("genreate string error:%w", err)
	}
	toString := base64.URLEncoding.EncodeToString(bs)
	return toString, nil
}

const SessionTokenBytes = 32

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
