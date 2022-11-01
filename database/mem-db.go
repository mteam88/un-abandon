package database

import (
	"fmt"
)

// MemDB is a simple in-memory database
type MemDB struct {
	data map[string][]byte
}

type User struct {
	Username string
	Token	string
	GithubID int64
}

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"html_url"`
	ID          int64  `json:"id"`
}

// NewMemDB creates a new MemDB
func NewMemDB() *MemDB {
	return &MemDB{
		data: make(map[string][]byte),
	}
}

// Get returns the value for a given key
func (db *MemDB) Get(key string) ([]byte, error) {
	val, ok := db.data[key]
	if !ok {
		return []byte{}, fmt.Errorf("key not found")
	}
	return val, nil
}

// Set sets the value for a given key
func (db *MemDB) Set(key string, val []byte) error {
	db.data[key] = val
	return nil
}

// Delete deletes the value for a given key
func (db *MemDB) Delete(key string) error {
	delete(db.data, key)
	return nil
}

// Close closes the database
func (db *MemDB) Close() error {
	return nil
}

func (db *MemDB) Dump() string {
	return fmt.Sprintf("%v", db.data)
}