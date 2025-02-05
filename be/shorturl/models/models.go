package models

import "time"

type (
	CreateShortURL struct {
		LongURL string `json:"long_url"`
	}

	Key struct {
		Key string `json:"key"`
	}

	DeleteShortURL struct {
		Key string `json:"key"`
	}

	UpdateShortURL struct {
		Key     string `json:"key"`
		LongURL string `json:"long_url"`
	}

	ShortURLRecord struct {
		Key       string    `json:"key"`
		LongURL   string    `json:"long_url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
