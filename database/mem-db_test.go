package database

// test mem-db.go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemDB(t *testing.T) {
	db := NewMemDB()

	err := db.Set("key", []byte{1,2,3})
	assert.Nil(t, err)

	val, err := db.Get("key")
	assert.Nil(t, err)
	assert.Equal(t, []byte{1,2,3}, val)

	err = db.Delete("key")
	assert.Nil(t, err)

	_, err = db.Get("key")
	assert.NotNil(t, err)

	err = db.Close()
	assert.Nil(t, err)
}