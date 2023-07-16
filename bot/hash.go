package bot

import (
	"crypto/md5"
	"encoding/hex"
)

func (b *Bot) textToHash(text string) (hash string) {
	h := md5.New()
	h.Write([]byte(text))
	hash = hex.EncodeToString(h.Sum(nil))
	b.hashCache[hash] = text
	return hash
}

func (b *Bot) hashToText(hash string) (text string) {
	return b.hashCache[hash]
}
