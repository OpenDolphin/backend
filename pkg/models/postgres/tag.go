package pg_model

type Tag struct {
	ID   uint64 `gorm:"primaryKey" json:"id"`
	Text string `gorm:"index,unique" json:"text"`

	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}
