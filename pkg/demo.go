package server

import (
	"fmt"
	pg_model "github.com/denysvitali/social/backend/pkg/models/postgres"
	"github.com/denysvitali/social/backend/pkg/unsplash"
	"github.com/oklog/ulid/v2"
)

func (s *Server) createDemoUsers() error {

	type inputUser struct {
		ID          uint64
		Username    string
		DisplayName string
		Bio         string
		ProfilePic  string
		Verified    bool
	}

	users := []inputUser{
		{
			ID:          1,
			Username:    "pdoe",
			DisplayName: "Pierre Doe",
			Bio:         "Just a random guy testing out this platform",
			ProfilePic:  unsplash.GetProfilePicture("0J9l9xRyMSo"),
			Verified:    true,
		},
		{
			ID:          2,
			Username:    "jfadel",
			DisplayName: "Johanna Fadel",
			Bio: "ust a small town girl living in a lonely world. Posting stuff about life, love, " +
				"and all the adventures in between.",
			ProfilePic: unsplash.GetProfilePicture("AJIqZDAUD7A"),
			Verified:   true,
		},
		{
			ID:          3,
			Username:    "sarahj",
			DisplayName: "Sarah Johnson",
			Bio:         "Writer, reader, and tea lover. Sharing my thoughts and adventures on this platform.",
			ProfilePic:  unsplash.GetProfilePicture("rDEOVtE7vOs"),
			Verified:    false,
		},
		{
			ID:          4,
			Username:    "mikem",
			DisplayName: "Mike Mitchell",
			Bio:         "Tech enthusiast, gamer, and movie buff. Sharing my love for all things geeky on this platform.",
			ProfilePic:  unsplash.GetProfilePicture("iFgRcqHznqg"),
			Verified:    true,
		},
		{
			ID:          5,
			Username:    "katieb",
			DisplayName: "Katie Brown",
			Bio:         "Traveler, foodie, and dog mom. Sharing my adventures and delicious finds on this platform.",
			ProfilePic:  unsplash.GetProfilePicture("O3ymvT7Wf9U"),
			Verified:    false,
		},
	}

	for _, u := range users {
		user := pg_model.User{
			ID:          u.ID,
			Username:    u.Username,
			DisplayName: u.DisplayName,
			Biography:   u.Bio,
			ProfilePictures: []pg_model.ProfilePicture{
				{
					Url: u.ProfilePic,
				},
			},
			Verified: u.Verified,
		}
		tx := s.pgDB.Create(&user)
		if tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}

func (s *Server) createDemoPosts() error {
	type post struct {
		ID      string
		Author  uint64
		LikedBy []uint64
		Content string
	}
	posts := []post{
		{
			ID:      "01GKF52B2HRVS4NWNY98M3D6RF",
			Author:  1,
			LikedBy: []uint64{1, 2, 5},
			Content: "It's good to be back!",
		},
		{
			ID:      "01GM3FZQB60JP3W1A8MV6TFHGZ",
			Author:  2,
			LikedBy: []uint64{1},
			Content: "About to enjoy the sunset 🌇",
		},
		{
			ID:      "01GM3GNC8RCC0B1C6HXXWB82W0",
			Author:  3,
			LikedBy: []uint64{2, 5},
			Content: "Anyone wanting to meet for a ski-trip?",
		},
	}

	for _, p := range posts {
		postUlid, err := ulid.Parse(p.ID)
		if err != nil {
			return fmt.Errorf("unable to parse post ulid: %v", err)
		}

		var likedBy []pg_model.User
		for _, id := range p.LikedBy {
			likedBy = append(likedBy, pg_model.User{ID: id})
		}

		tx := s.pgDB.Create(&pg_model.Post{
			ID:       postUlid.Bytes(),
			AuthorID: p.Author,
			Content:  p.Content,
			LikedBy:  likedBy,
		})

		if tx.Error != nil {
			return err
		}
	}
	return nil
}
