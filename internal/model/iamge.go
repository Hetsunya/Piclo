package model

import "time"

type Image struct {
	ID         string
	Filename   string
	MIMEType   string
	Size       int64
	StorageKey string
	CreatedAt  time.Time
	ExpiresAt  *time.Time
}
