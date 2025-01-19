package service

import (
	"database/sql"
	"errors"
	"html"
	"strconv"
	"strings"

	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	utils "real-time-forum/pkg"
)

type Service struct {
	Database *repository.Database
}

func (s *Service) CreatePost(post models.Post, uid string) error {
	var err error
	post.UserID, err = s.Database.GetUser(uid)
	if err != nil {
		return err
	}

	post.Categories = removeDuplicate(post.Categories)
	err = CheckPostValidation(post)
	if err != nil {
		return err
	}

	post.Title = html.EscapeString(post.Title)
	post.Content = html.EscapeString(post.Content)

	postId, err := s.Database.InsertPost(post)
	if err != nil {
		return err
	}

	err = s.Database.AddCategoriesToPost(postId, post.Categories)
	if err != nil {
		if errDB := s.Database.DeletePost(postId); errDB != nil {
			return errDB
		}
		return err
	}
	return nil
}

func CheckPostValidation(post models.Post) error {
	if len(strings.TrimSpace(post.Title)) == 0 || len(post.Title) > 500 {
		return errors.New("models.PostErrors.TitleLength")
	}
	if len(strings.TrimSpace(post.Content)) == 0 || len(post.Content) > 5000 {
		return errors.New("models.PostErrors.ContentLength")
	}
	if post.UserID == 0 {
		return errors.New("models.UserErrors.UserNotExist")
	}
	return nil
}

func removeDuplicate(categories []string) []string {
	var result []string
	for _, categorie := range categories {
		if !utils.Contains(result, categorie) {
			result = append(result, categorie)
		}
	}
	return result
}

func (s *Service) GetPostbyid(idstr string, userid int) (models.Post, error) {
	postid, err := strconv.Atoi(idstr)
	if err != nil {
		return models.Post{}, errors.New("not found")
	}

	post, err := s.Database.GetPost(postid, userid)
	if err != nil {
		return models.Post{}, err
	}

	categories, err := s.Database.GetPostCategories(postid)
	if err != nil {
		return models.Post{}, err
	}

	post.Categories = categories

	return post, nil
}

func (s *Service) GetPost(num, userID int) ([]models.Post, error) {
	start := (num * models.PostsPerPage)
	total := 0
	err := s.Database.Tablelen("post", &total)
	if err != nil {
		return nil, err
	}
	if start > total {
		return nil, sql.ErrNoRows
	}
	row, err := s.Database.ExtractPosts(start)
	if err != nil {
		return nil, err
	}
	var posts []models.Post
	for row.Next() {
		var post models.Post

		err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Author, &post.Likes, &post.Dislikes, &post.CommentsCount)
		if err != nil {
			return nil, err
		}
		post.IsLiked, post.IsDisliked = s.Database.CheckIfLikedPost(post.ID, userID)
		// Get categories
		categories, err := s.Database.GetPostCategories(post.ID)
		if err != nil {
			return nil, err
		}

		post.Categories = categories

		posts = append(posts, post)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}