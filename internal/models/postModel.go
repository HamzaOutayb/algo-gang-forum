package models

import "time"

type Post struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	Created       time.Time `json:"date"`
	Likes         int       `json:"likes"`
	Dislikes      int       `json:"dislikes"`
	IsLiked       bool      `json:"isliked"`
	IsDisliked    bool      `json:"isdisliked"`
	CommentsCount int       `json:"commentsCount"`
	Categories    []string  `json:"categories"`
	Joined_at     time.Time `json:"joined_at"`
}

var PostErrors struct {
	PostNotExist string
	ContentLength string
	TitleLength string
	CategoryDoesntExist string
} = struct{ PostNotExist string; ContentLength string; TitleLength string; CategoryDoesntExist string }{
	PostNotExist: "post doesn't exist",
	ContentLength: "invalid content length",
	TitleLength: "invalid title length",
	CategoryDoesntExist: "categories doesn't exist",
}