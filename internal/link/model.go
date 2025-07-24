package link

import (
	"math/rand"

	"gorm.io/gorm"
)

type Link struct {
	*gorm.Model
	Url  string `json:"url"`
	Hash string `json:"hash" gorm:"uniqeIndex"`
}

func NewLink(url string) *Link {
	return &Link{
		Url:  url,
		Hash: RandStringRunes(5),
	}
}

var letterRunes = []rune("qqwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Int()]
	}
	return string(b)
}
