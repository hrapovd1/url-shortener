package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hrapovd1/url-shortener/internal/app/errors"
)

func TestNewMemStorage(t *testing.T) {
	stor := NewMemStorage()
	require.IsType(t, &MemStorage{}, stor)
}

func TestMemStorage_GetShort(t *testing.T) {
	stor := NewMemStorage()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"check short generator", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stor.GetShort()
			require.NoError(t, err)
			assert.Equal(t, strLen, len([]rune(got)))
		})
	}
}

func TestMemStorage_SaveURL(t *testing.T) {
	short := "abcd"
	url := "ya.ru"
	stor := NewMemStorage()
	require.NoError(t, stor.SaveURL(url, short))
	urlResult, ok := map[string]string(*stor)[short]
	require.True(t, ok)
	assert.Equal(t, url, urlResult)
}

func TestMemStorage_GetURL(t *testing.T) {
	stor := NewMemStorage()
	require.NoError(t, stor.SaveURL("ya.ru", "abcd"))

	tests := []struct {
		name    string
		short   string
		want    string
		errored bool
	}{
		{"exist url", "abcd", "ya.ru", false},
		{"not exist url", "bcdef", "", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url, err := stor.GetURL(test.short)
			if test.errored {
				assert.ErrorIs(t, err, errors.ErrorStorageGetShort)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.want, url)
		})
	}
}

func TestRandSeq(t *testing.T) {
	result := randSeq(strLen)
	assert.IsType(t, "", result)
	assert.Len(t, []rune(result), strLen)
}
