package pg_model

type Post struct {
	// ID is an ULID that contains the post creation date and some randomness
	ID      []byte `gorm:"primaryKey,type:bytea" json:"id"`
	Content string `json:"content"`

	ParentPost   *Post  `json:"parentPost,omitempty"`
	ParentPostID []byte `json:"parentPostID,type:bytea,omitempty"`

	UserMention []User `gorm:"many2many:user_mention;" json:"userMention"`
	Tags        []Tag  `gorm:"many2many:post_tags;" json:"tags"`
	LikedBy     []User `gorm:"many2many:user_likes" json:"liked_by,omitempty"`

	Likes uint64 `json:"likes"`

	Author   *User  `json:"author,omitempty" gorm:"foreignkey:AuthorID"`
	AuthorID uint64 `json:"authorId"`
	Deleted  bool   `json:"-"`
}
