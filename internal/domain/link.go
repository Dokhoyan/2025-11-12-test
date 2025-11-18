package domain

import "time"

type LinkSet struct {
	ID        int64     `json:"id"`
	Links     []Link    `json:"links"`
	CreatedAt time.Time `json:"created_at"`
}

type Link struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

type LinkStatus string

const (
	StatusAvailable   LinkStatus = "available"
	StatusUnavailable LinkStatus = "unavailable"
	StatusChecking    LinkStatus = "checking"
)
