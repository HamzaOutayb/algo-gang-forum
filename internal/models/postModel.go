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

type ReactResponse struct {
	Like       int  `json:"Like"`
	Dislike    int  `json:"Dislike"`
	IsLiked    bool `json:"isliked"`
	IsDisliked bool `json:"isdisliked"`
}

type Comment struct {
	UserId  int
	UserUID string
	PostId  int    `json:"postId"`
	Content string `json:"content"`
}

type ShowComment struct {
	Id         int    `json:"id"`
	Author     string `json:"author"`
	Content    string `json:"content"`
	Likes      int    `json:"likes"`
	Dislikes   int    `json:"dislikes"`
	Date       string `json:"date"`
	IsLiked    bool   `json:"isliked"`
	IsDisliked bool   `json:"isdisliked"`
}

var PostsPerPage = 5

var PostErrors struct {
	PostNotExist        string
	ContentLength       string
	TitleLength         string
	CategoryDoesntExist string
} = struct {
	PostNotExist        string
	ContentLength       string
	TitleLength         string
	CategoryDoesntExist string
}{
	PostNotExist:        "post doesn't exist",
	ContentLength:       "invalid content length",
	TitleLength:         "invalid title length",
	CategoryDoesntExist: "categories doesn't exist",
}

var CommentErrors struct {
	InvalidCommentLength string
	InvalidPage          string
} = struct {
	InvalidCommentLength string
	InvalidPage          string
}{
	InvalidCommentLength: "invalid comment length",
	InvalidPage:          "invalid page number",
}
