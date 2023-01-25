package pg_model

type Post struct {
	// ID is an ULID that contains the post creation date and some randomness
	ID      []byte `gorm:"primaryKey,type:bytea" json:"id"`
	Content string `json:"content"`

	UserMention []User `gorm:"many2many:user_mention;" json:"userMention"`
	Tags        []Tag  `gorm:"many2many:post_tags;" json:"tags"`

	// Do not rely on these intensive operations, unless we really need to
	LikedBy []User `gorm:"many2many:user_likes" json:"-,omitempty"`

	ParentPostID *[]byte
	ChildPosts   []Post `gorm:"foreignKey:ParentPostID"`

	Likes    uint64 `json:"likes"`
	Reshares uint64 `json:"reshares"`

	Author   *User  `json:"author,omitempty" gorm:"foreignkey:AuthorID"`
	AuthorID uint64 `json:"authorId"`
	Deleted  bool   `json:"-"`
}
