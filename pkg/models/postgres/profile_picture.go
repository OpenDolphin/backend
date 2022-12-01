package pg_model

import "time"

type ProfilePicture struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserId      uint       `json:"user_id"`
	LastUpdated *time.Time `json:"last_updated,omitempty"`
	Url         string     `json:"url"`
}
