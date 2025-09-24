package link

import (
	"go-advanced/internal/stat"
	"math/rand"

	"gorm.io/gorm"
)

type Link struct {
	*gorm.Model
	Url   string      `json:"url"`
	Hash  string      `json:"hash" gorm:"uniqeIndex"`
	Stats []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewLink(url string) *Link {
	link := &Link{
		Url: url,
	}
	link.GenerateHash()
	return link
}

func (link *Link) GenerateHash() {
	link.Hash = RandStringRunes(5)
}

var letterRunes = []rune("qqwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
