package database

import (
	"encoding/json"
	"log"
)

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
	Token		string `json:"token"`
}

// define json serializer for Repo
func (r Repo) Serialize() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Print(err)
	}
	return b
}

// define json deserializer for Repo
func DeserializeRepo(b []byte) Repo {
	var r Repo
	err := json.Unmarshal(b, &r)
	if err != nil {
		log.Print(err)
	}
	return r
}

// get bytes id from a Github repository ID
func GetID(ghid int32) []byte {
	return []byte(string(ghid))
}

func (u User) Serialize() []byte {
	b, err := json.Marshal(u)
	if err != nil {
		log.Print(err)
	}
	return b
}

func DeserializeUser(b []byte) User {
	var u User
	err := json.Unmarshal(b, &u)
	if err != nil {
		log.Print(err)
	}
	return u
}