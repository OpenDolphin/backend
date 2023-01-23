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
		BioPic      string
		ProfilePic  string
		Verified    bool
	}

	users := []inputUser{
		{
			ID:          1,
			Username:    "jdoe",
			DisplayName: "John Doe",
			Bio:         "Just a random guy testing out this platform",
			BioPic:      unsplash.GetBioPic("ukvgqriuOgo"),
			ProfilePic:  unsplash.GetProfilePicture("lETfyhB8g4Q"),
			Verified:    true,
		},
		{
			ID:          2,
			Username:    "sjohnson",
			DisplayName: "Samantha Johnson",
			Bio:         "Marketing professional and social media enthusiast",
			BioPic:      unsplash.GetBioPic("Nyvq2juw4_o"),
			ProfilePic:  unsplash.GetProfilePicture("O3ymvT7Wf9U"),
			Verified:    true,
		},
		{
			ID:          3,
			Username:    "mwilson",
			DisplayName: "Mike Wilson",
			Bio:         "Professional photographer and travel blogger",
			BioPic:      unsplash.GetBioPic("LY1eyQMFeyo"),
			ProfilePic:  unsplash.GetProfilePicture("hh3ViD0r0Rc"),
			Verified:    false,
		},
		{
			ID:          4,
			Username:    "paula_g",
			DisplayName: "Paula Garcia",
			Bio:         "Entrepreneur and business owner",
			BioPic:      unsplash.GetBioPic("99SXcea3uOk"),
			ProfilePic:  unsplash.GetProfilePicture("cUKy1J3wzqg"),
			Verified:    true,
		},
		{
			ID:          5,
			Username:    "tjones",
			DisplayName: "Tina Jones",
			Bio:         "Freelance writer and book lover",
			BioPic:      unsplash.GetBioPic("f4845LpnSbs"),
			ProfilePic:  unsplash.GetProfilePicture("zNWlX5Sw9a4"),
			Verified:    false,
		},
		{
			ID:          6,
			Username:    "samantha_j",
			DisplayName: "Samantha Johnson",
			Bio:         "Marketing professional and social media enthusiast",
			BioPic:      unsplash.GetBioPic("KEHVSsRtnL0"),
			ProfilePic:  unsplash.GetProfilePicture("-zqoE7jnQgw"),
			Verified:    true,
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
			BiographyPictures: []pg_model.BioPicture{
				{
					Url: u.BioPic,
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
			ID:      "01GQGDW0E19PRG21S0DWS090X0",
			Author:  1,
			LikedBy: []uint64{2, 3, 4, 5},
			Content: "Just finished a great hike in the mountains! #nature #adventure",
		},
		{
			ID:      "01GQGDW16BGX8K6Q4NGNNTSNAE",
			Author:  2,
			LikedBy: []uint64{1, 3},
			Content: "Can't wait for the weekend! Anyone up for a road trip? #roadtrip #friends",
		},
		{
			ID:      "01GQGDW8F05WF7RQ8ZKTZM6TAW",
			Author:  3,
			LikedBy: []uint64{1},
			Content: "Just had the best sushi of my life! #foodie #yum",
		},
		{
			ID:      "01GQGDWCVBXJK7Y16DNZ5M0B44",
			Author:  4,
			LikedBy: []uint64{},
			Content: "Who's ready for some live music tonight? #concert #music",
		},
		{
			ID:      "01GQGDWMVDFMKED8Z7T88PEY6Q",
			Author:  5,
			LikedBy: []uint64{2, 4},
			Content: "Just landed in a new city! Excited to explore! #travel #citylife",
		},
		{
			ID:      "01GQGDWZSR0K3B6NMFHM2WTDBE",
			Author:  1,
			LikedBy: []uint64{2},
			Content: "Who's ready for some football? #sports #gameon",
		},
		{
			ID:      "01GQGDXA5FT6CND4YYJJYB5HYP",
			Author:  2,
			LikedBy: []uint64{},
			Content: "Just finished reading a great book! #reading #bookworm",
		},
		{
			ID:      "01GQGDXR7J8WJTNG4E6QAP5T02",
			Author:  3,
			LikedBy: []uint64{6},
			Content: "Who's up for a round of golf this weekend? #golf #weekendfun",
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
