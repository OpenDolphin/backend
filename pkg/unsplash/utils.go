package unsplash

import (
	"fmt"
	"net/url"
)

func GetProfilePicture(id string) string {
	return fmt.Sprintf("https://unsplash.com/photos/%s/download?w=200", url.PathEscape(id))
}

func GetBioPic(id string) string {
	return fmt.Sprintf("https://unsplash.com/photos/%s/download?w=800", url.PathEscape(id))
}
