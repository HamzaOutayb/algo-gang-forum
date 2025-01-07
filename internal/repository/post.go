package repository

import (
	"errors"

	"real-time-forum/internal/models"
)

func (database *Database) InsertPost(post models.Post) (int, error) {
	rowrResult, err := database.Db.Exec("INSERT INTO post (title, content, user_id) VALUES (?, ?, ?)", post.UserID, post.Title, post.Content, post.Joined_at)
	if err != nil {
		return 0, err
	}
	postid, err := rowrResult.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(postid), nil
}

func (database *Database) AddCategoriesToPost(Postid int, categories []string) error {
	for _, categorie := range categories {
		categorieId, err := database.GetCategoryId(categorie)
		if err != nil {
			return err
		}
		if categorieId == 0 {
			return errors.New(models.PostErrors.CategoryDoesntExist)
		}
		// add the category to the post using category_post table
		err = database.AddCategory(Postid, categorieId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) AddCategory(postId, categoryId int) error {
	_, err := d.Db.Exec("INSERT INTO post_category (post_id, category_id) VALUES(?, ?)", postId, categoryId)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetCategoryId(name string) (int, error) {
	var id int
	err := d.Db.QueryRow("SELECT id FROM categories WHERE category_name = ?", name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *Database) DeletePost(post_id int) error {
	row, err := d.Db.Prepare("DELETE FROM post WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = row.Exec(post_id)
	if err != nil {
		return err
	}
	return nil
}
