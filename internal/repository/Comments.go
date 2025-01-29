package repository

import "real-time-forum/internal/models"

const CommentsPerPage = 15;

// Get comments count
func (database *Database) GetCommentsCount(postId int) (int, error) {
	var count int
	err := database.Db.QueryRow("SELECT COUNT(*) FROM comment WHERE post_id = ?", postId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Get Comments from a specific comment row (like from 1 and get the 100 in front of it)
func (database *Database) GetCommentsFrom(from, postId, userId int) ([]models.ShowComment, error) {
	rows, err := database.Db.Query(
		`SELECT comment_id, comment_author, comment_content, comment_date, comment_likes, comment_dislikes
		FROM single_comment
		WHERE post_id = ? ORDER BY comment_date DESC LIMIT ? OFFSET ?`,
		postId, CommentsPerPage, from)
	if err != nil {
		return nil, err
	}

	var comments []models.ShowComment
	for rows.Next() {
		var comment models.ShowComment
		rows.Scan(&comment.Id, &comment.Author, &comment.Content, &comment.Date, &comment.Likes, &comment.Dislikes)
		comment.IsLiked, comment.IsDisliked = database.CheckIfLikedComment(comment.Id, userId)
		comments = append(comments, comment)
	}

	return comments, err
}