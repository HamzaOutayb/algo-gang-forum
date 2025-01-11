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
}

type React struct {
	Thread_type string `json:"thread_type"`
	Thread_id   int    `json:"thread_id"`
	React       int    `json:"react"`
}

var PostsPerPage = 20

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