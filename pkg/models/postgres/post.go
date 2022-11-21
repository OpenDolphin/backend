package pg_model

import "time"

type Post struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`

	ParentPost   *Post  `json:"parentPost,omitempty"`
	ParentPostID uint64 `json:"parentPostID,omitempty"`

	ReshareCount uint32 `json:"reshareCount"`
	ReplyCount   uint32 `json:"replyCount"`

	UserMention []User `gorm:"many2many:user_mention;" json:"userMention"`
	Tags        []Tag  `gorm:"many2many:post_tags;" json:"tags"`

	Author   *User  `json:"author,omitempty"`
	AuthorID uint64 `json:"authorId"`
	Deleted  bool   `json:"-"`
}
