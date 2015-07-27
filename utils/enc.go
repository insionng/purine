package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func Md5String(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

func Sha256String(str string) string {
	s := sha256.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

func Base64EncodeBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	defer buf.Reset()
	return buf.Bytes(), nil
}
