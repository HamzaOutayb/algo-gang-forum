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

func (d *Database) GetPost(id, user_id int) (models.Post, error) {
	var post models.Post

	row := d.Db.QueryRow(`SELECT  post_id, post_title, post_content, post_date, post_author, post_likes, post_dislikes, post_comments_count, joined_at
	FROM single_post
	WHERE post_id = ?`, id)

	if row.Err() != nil {
		return models.Post{}, row.Err()
	}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Author, &post.Likes, &post.Dislikes, &post.CommentsCount, &post.Joined_at)
	if err != nil {
		return models.Post{}, row.Err()
	}
	if id != 0 {
		post.IsLiked, post.IsDisliked = d.CheckIfLikedPost(post.ID, user_id)
	}
	return post, row.Err()
}

func (d *Database) CheckIfLikedPost(post_id, user_id int) (bool, bool) {
	isLiked := 0
	d.Db.QueryRow("SELECT is_liked FROM postReact WHERE post_id = ? AND user_id = ?", post_id, user_id).Scan(&isLiked)
	switch isLiked {
	case 2:
		return false, true
	case 1:
		return true, false
	default:
		return false, false
	}
}

func (a *Database) GetPostCategories(postId int) ([]string, error) {
	// Get Categories Ids
	var names []string
	rows, err := a.Db.Query(`
	SELECT 
    c.category_name
FROM 
    post p
JOIN 
    post_category pc ON pc.post_id = p.id
JOIN
    categories c ON c.id = pc.category_id
WHERE p.id = ?`, postId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	return names, nil
}