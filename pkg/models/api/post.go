package api

import "time"

type Post struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Likes     uint64    `json:"likes"`
	Replies   []Post    `json:"replies"`
	Author    uint64    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

type PostsResponse struct {
	Posts []Post `json:"posts"`
	Users []User `json:"users"`
}
