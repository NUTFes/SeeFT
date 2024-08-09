package entity

import (
	"time"
)

type Token struct {
	AccessToken string `json:"accessToken"`
}

type IsSignIn struct {
	IsSignIn bool `json:"isSignIn"`
}

type Session struct {
	ID          int
	AuthID      int
	UserID      int
	AccessToken string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}


