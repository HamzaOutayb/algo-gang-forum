package service

import (
	"errors"
	"html"
	"strings"

	"real-time-forum/internal/models"
)

const CommentsPerPage = 15

func (s *Service) GetComments(postId, page, userId int) ([]models.ShowComment, error) {
	// Validate page number
	if page < 0 {
		return nil, errors.New(models.CommentErrors.InvalidPage)
	}

	// Transfer "page" to "from" (page 1 mean page one that has 100 comment from 1 mean comment 1)
	from := (CommentsPerPage * page) - CommentsPerPage

	// Get the comments count number to check if the page number is right
	commentsCount, err := s.Database.GetCommentsCount(postId)
	if err != nil {
		return nil, err
	}

	if page-1 == (commentsCount/models.PostsPerPage)+(commentsCount%models.PostsPerPage) {
		from = models.PostsPerPage % commentsCount
	} else if page-1 > (commentsCount/models.PostsPerPage)+(commentsCount%models.PostsPerPage) {
		return []models.ShowComment{}, nil
	}
	// Get comments
	return s.Database.GetCommentsFrom(from, postId, userId)
}

func (s *Service) AddComment(comment models.Comment) error {
	// add the userId to the comment
	userID, _ := s.Database.GetUser(comment.UserUID)
	comment.UserId = userID

	// check if the post exist using the CheckPostExist
	if !s.Database.CheckPostExist(comment.PostId) {
		return errors.New(models.PostErrors.PostNotExist)
	}

	// Trim the space from the comment content
	comment.Content = strings.TrimSpace(comment.Content)

	// Fix html
	comment.Content = html.EscapeString(comment.Content)

	// Add the comment Using InsertComment
	err := s.Database.InsertComment(comment)

	return err
}
