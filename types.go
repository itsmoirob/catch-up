package main

import (
	"time"
)

type UpdateMediaRequest struct {
	Title     string `json:"title"`
	Form      string `json:"form"`
	IsWatched bool   `json:"isWatched"`
}

type CreateMediaRequest struct {
	Title string `json:"title"`
	Form  string `json:"form"`
}

type Media struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Form        string    `json:"form"`
	IsWatched   bool      `json:"isWatched"`
	DateWatched time.Time `json:"dateWatched"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewMedia(title, form string) *Media {
	return &Media{
		Title:     title,
		Form:      form,
		IsWatched: false,
		CreatedAt: time.Now().UTC(),
	}
}
