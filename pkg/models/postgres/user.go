package pg_model

import "time"

type User struct {
	ID              uint64           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time        `json:"createdAt"`
	Username        string           `gorm:"index" json:"username"`
	DisplayName     string           `json:"displayName"`
	Biography       string           `json:"biography"`
	Location        string           `json:"location"`
	FollowersCount  int              `json:"followersCount"`
	FollowingCount  int              `json:"followingCount"`
	ProfilePictures []ProfilePicture `json:"profilePictures,omitempty"`
	Verified        bool             `json:"verified"`
	Deleted         bool             `json:"-"`

	MentionedIn []Post `gorm:"many2many:user_mention;" json:"mentionedIn"`

	// HasMany relations
	Posts []Post `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
}
