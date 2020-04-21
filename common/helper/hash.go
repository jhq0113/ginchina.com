package helper

import (
	"bytes"
	"crypto"
	"encoding/hex"
)

func Hash(hash crypto.Hash, data []byte) []byte {
	h := hash.New()
	h.Write(data)
	return h.Sum(nil)
}

func HashString(hash crypto.Hash, data string) string {
	return HashAndHexEncode(hash, bytes.NewBufferString(data).Bytes())
}

func HashAndHexEncode(hash crypto.Hash, data []byte) string {
	return hex.EncodeToString(Hash(hash, data))
}

func Md5(data string) string {
	return HashString(crypto.MD5, data)
}

func Sha1(data string) string {
	return HashString(crypto.SHA1, data)
}

func Sha256(data string) string {
	return HashString(crypto.SHA256, data)
}

func Sha512(data string) string {
	return HashString(crypto.SHA512, data)
}
