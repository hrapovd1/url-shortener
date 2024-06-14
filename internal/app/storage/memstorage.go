package storage

import (
	"math/rand"

	"github.com/hrapovd1/url-shortener/internal/app/errors"
)

type MemStorage map[string]string

func NewMemStorage() *MemStorage {
	memStorage := make(MemStorage)
	return &memStorage
}

func (ms *MemStorage) GetShort() (string, error) {
	var short string
	for i := 0; i < 3; i++ {
		shortTmp := randSeq(strLen)
		if _, exist := map[string]string(*ms)[shortTmp]; exist {
			continue
		}
		short = shortTmp
	}
	if short == "" {
		return short, errors.ErrorStorageGenShort
	}
	return short, nil
}

func (ms *MemStorage) GetURL(short string) (string, error) {
	url, ok := map[string]string(*ms)[short]
	if !ok {
		return url, errors.ErrorStorageGetShort
	}
	return url, nil
}

func (ms *MemStorage) SaveURL(url string, short string) error {
	map[string]string(*ms)[short] = url
	return nil
}

const strLen = 9

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
