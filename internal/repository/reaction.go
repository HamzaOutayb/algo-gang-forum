package repository

import (
	"real-time-forum/internal/models"
)

func (d *Database) CheckIfThreadExists(react models.React) bool {
	exist := false
	if react.Thread_type == "post" {
		d.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM post WHERE id = ?)", react.Thread_id).Scan(&exist)
		return exist
	} else {
		d.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM comment WHERE id = ?)", react.Thread_id).Scan(&exist)
		return exist
	}
}

func (d *Database) ReactionService(react models.React, user_id int) error {
	if react.Thread_type == "post" {
		err := d.postReaction(react.Thread_id, user_id, react.React)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := d.commentReaction(react.Thread_id, user_id, react.React)
		if err != nil {
			return err
		}
		return nil
	}
}

func (data *Database) commentReaction(comment_id, user_id, react int) error {
	var exists bool
	exists, err := data.CheckCommentReaction(user_id, comment_id)
	if err != nil {
		return err
	}
	if !exists {
		err := data.InsertReactComment(user_id, comment_id, react)
		if err != nil {
			return err
		}
	} else {
		var isLiked int
		isLiked, err := data.GetReactionTypeComment(user_id, comment_id)
		if err != nil {
			return err
		}
		if isLiked == react {
			err := data.DeleteReactionComment(user_id, comment_id)
			if err != nil {
				return err
			}
		} else {
			err := data.DeleteReactionComment(user_id, comment_id)
			if err != nil {
				return err
			}
			err = data.InsertReactComment(user_id, comment_id, react)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (data *Database) postReaction(post_id, user_id, react int) error {
	var exists bool
	exists, err := data.CheckPostReaction(user_id, post_id)
	if err != nil {
		return err
	}
	if !exists {
		err := data.InsertReactPost(user_id, post_id, react)
		if err != nil {
			return err
		}
	} else {
		var like_type int
		like_type, err := data.GetReactionTypePost(user_id, post_id)
		if err != nil {
			return err
		}
		if like_type == react {
			err := data.DeleteReactionPost(user_id, post_id)
			if err != nil {
				return err
			}
		} else {
			err := data.DeleteReactionPost(user_id, post_id)
			if err != nil {
				return err
			}
			err = data.InsertReactPost(user_id, post_id, react)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (db *Database) CheckPostReaction(user_id, post_id int) (bool, error) {
	var exists bool
	err := db.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM postReact WHERE user_Id = ? AND post_Id = ?)", user_id, post_id).Scan(&exists)
	return exists, err
}

func (db *Database) CheckCommentReaction(user_id, comment_id int) (bool, error) {
	var exists bool
	err := db.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM commentReact WHERE user_Id = ? AND comment_Id = ?)", user_id, comment_id).Scan(&exists)
	return exists, err
}

func (db *Database) DeleteReactionPost(user_id, post_id int) error {
	_, err := db.Db.Exec("DELETE FROM postReact WHERE post_Id = ? AND user_Id = ?", post_id, user_id)
	return err
}

func (db *Database) DeleteReactionComment(user_id, post_id int) error {
	_, err := db.Db.Exec("DELETE FROM commentReact WHERE comment_Id = ? AND user_Id = ?", post_id, user_id)
	return err
}

func (db *Database) GetReactionTypePost(user_id, post_id int) (int, error) {
	var isLiked int
	err := db.Db.QueryRow("SELECT is_liked FROM postReact WHERE user_id = ? AND post_Id = ?", user_id, post_id).Scan(&isLiked)
	return isLiked, err
}

func (db *Database) GetReactionTypeComment(user_id, post_id int) (int, error) {
	var isLiked int
	err := db.Db.QueryRow("SELECT is_liked FROM commentReact WHERE user_id = ? AND comment_Id = ?", user_id, post_id).Scan(&isLiked)
	return isLiked, err
}

func (db *Database) InsertReactPost(user_id, post_id, like_type int) error {
	_, err := db.Db.Exec("INSERT INTO postReact (post_Id, user_Id, is_liked) VALUES (?,?,?)", post_id, user_id, like_type)
	return err
}

func (db *Database) InsertReactComment(user_id, post_id, like_type int) error {
	_, err := db.Db.Exec("INSERT INTO commentReact (comment_Id, user_Id, is_liked) VALUES (?,?,?)", post_id, user_id, like_type)
	return err
}

func (db *Database) CountPostLikes(post_id int) (int, int, error) {
	var likes, dislikes int
	err := db.Db.QueryRow("SELECT COUNT(*) FROM postReact WHERE is_liked = 1 AND post_Id = ?", post_id).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}

	err = db.Db.QueryRow("SELECT COUNT(*) FROM postReact WHERE is_liked = 2 AND post_Id = ?", post_id).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

func (db *Database) CountCommentLikes(post_id int) (int, int, error) {
	var likes, dislikes int
	err := db.Db.QueryRow("SELECT COUNT(*) FROM commentReact WHERE is_liked = 1 AND comment_Id = ?", post_id).Scan(&likes)
	if err != nil {
		return 0, 0, err
	}

	err = db.Db.QueryRow("SELECT COUNT(*) FROM commentReact WHERE is_liked = 2 AND comment_Id = ?", post_id).Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}

	return likes, dislikes, nil
}

func (data *Database) LikesTotal(thread_type string, thread_id int) (models.ReactResponse, error) {
	response := models.ReactResponse{}
	var err error
	if thread_type == "post" {
		response.Like, response.Dislike, err = data.CountPostLikes(thread_id)
	} else {
		response.Like, response.Dislike, err = data.CountCommentLikes(thread_id)
	}
	if err != nil {
		return response, err
	}
	return response, nil
}

func (d *Database) CheckIfLikedComment(post_id, user_id int) (bool, bool) {
	isLiked := 0
	d.Db.QueryRow("SELECT is_liked FROM commentReact WHERE comment_id = ? AND user_id = ?", post_id, user_id).Scan(&isLiked)
	switch isLiked {
	case 2:
		return false, true
	case 1:
		return true, false
	default:
		return false, false
	}
}
