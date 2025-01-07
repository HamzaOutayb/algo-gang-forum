package service

import (
	"fmt"
	"html"
	"strings"

	"real-time-forum/internal/models"
	"real-time-forum/internal/repository"
	utils "real-time-forum/pkg"
)

type Service struct {
	Database *repository.Database
}

func (d *Service) CreatePost(post models.Post, uid string) error {
	var err error
	post.UserID, err = d.Database.GetUser(uid)
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

	postId, err := d.Database.InsertPost(post)
	if err != nil {
		return err
	}

	err = d.Database.AddCategoriesToPost(postId, post.Categories)
	if err != nil {
		if errDB := d.Database.DeletePost(postId); errDB != nil {
			return errDB
		}
		return err
	}
	return nil
}

func CheckPostValidation(post models.Post) error {
	if len(strings.TrimSpace(post.Title)) == 0 || len(post.Title) > 500 {
		return fmt.Errorf(models.PostErrors.TitleLength)
	}
	if len(strings.TrimSpace(post.Content)) == 0 || len(post.Content) > 5000 {
		return fmt.Errorf(models.PostErrors.ContentLength)
	}
	if post.UserID == 0 {
		return fmt.Errorf(models.UserErrors.UserNotExist)
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