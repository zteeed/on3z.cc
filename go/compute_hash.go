package app

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"
)

func removeCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if strings.IndexRune(characters, r) < 0 {
			return r
		}
		return -1
	}
	return strings.Map(filter, input)
}

func computeShortURL(LongURL string) string {
	hash := sha1.New()
	hash.Write([]byte(LongURL))
	hashBase64 := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	hashBase64Stripped := removeCharacters(hashBase64, "+/=")
	return hashBase64Stripped[:7]
}
