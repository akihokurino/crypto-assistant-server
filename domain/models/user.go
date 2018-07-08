package models

import "net/url"

type UserID string

type User struct {
	Id       UserID
	Username string
	IconURL  *url.URL
}

func NewUser(id UserID, username string) *User {
	return &User{
		Id:       id,
		Username: username,
	}
}
