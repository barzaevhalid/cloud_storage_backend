package models

import "time"

type File struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	Filename     string     `json:"filename"`
	OriginalName string     `json:"original_name"`
	MimeType     string     `json:"mime_type"`
	Size         int64      `json:"size"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
