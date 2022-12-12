package api

type User struct {
	ID          uint64 `json:"id"`
	DisplayName string `json:"displayName"`
	Username    string `json:"username"`
	Verified    bool   `json:"verified"`
}
