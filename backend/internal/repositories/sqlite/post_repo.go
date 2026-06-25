package sqlite

import (
	"database/sql"

	"social-network/backend/internal/models"
)

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) CreatePost(post *models.Post) error {
	_, err := r.db.Exec(`
		INSERT INTO posts (id, user_id, content, media_path, media_type, privacy, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, post.ID, post.UserID, post.Content, post.MediaPath, post.MediaType, post.Privacy, post.CreatedAt, post.UpdatedAt)
	return err
}

func (r *postRepository) CreateAllowedUser(postID, userID string) error {

	_, err := r.db.Exec(`
		INSERT INTO post_allowed_users (post_id, user_id)
		VALUES (?, ?)
	`, postID, userID)
	return err

}

func (r *postRepository) GetPostByID(postID string) (*models.Post, error) {
	row := r.db.QueryRow(`
    	SELECT id, user_id, content, media_path, media_type, privacy, created_at, updated_at
    	FROM posts
    	WHERE id = ?
	`, postID)

	var p models.Post

	err := row.Scan(&p.ID, &p.UserID, &p.Content, &p.MediaPath, &p.MediaType, &p.Privacy, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *postRepository) IsAllowedToViewPost(postID, userID string) (bool, error) {

	var count int

	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM post_allowed_users
		WHERE post_id = ? AND user_id = ?
	`, postID, userID).Scan(&count)

	if err != nil {
		return false, err
	}
	return count > 0, nil

}

func (r *postRepository) DeletePost(postID string) error {

	_, err := r.db.Exec(`
		DELETE FROM posts WHERE id = ?
	 `, postID)
	return err

}
func (r *postRepository) UpdatePost(post *models.Post) error {

	_, err := r.db.Exec(`
		UPDATE 	posts
		SET content = ?, media_path = ?, media_type = ?, privacy = ?, updated_at = ?
    	WHERE id = ?
	`, post.Content, post.MediaPath, post.MediaType, post.Privacy, post.UpdatedAt, post.ID)
	return err

}
func (r *postRepository) GetPostsByUserID(userID, viewerID string) ([]*models.Post, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, content, media_path, media_type, privacy, created_at, updated_at
		FROM posts
		WHERE user_id = ?
		AND (
			user_id = ?
			OR privacy = 'public'
			OR (privacy = 'followers' AND EXISTS (
				SELECT 1 FROM followers
				WHERE follower_id = ? AND following_id = ?
			))
			OR (privacy = 'selected' AND EXISTS (
				SELECT 1 FROM post_allowed_users
				WHERE post_id = posts.id AND user_id = ?
			))
		)
		ORDER BY created_at DESC
	`, userID, viewerID, viewerID, userID, viewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.MediaPath, &p.MediaType, &p.Privacy, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

// GetFeed now joins users to include author name and avatar
func (r *postRepository) GetFeed(userID string) ([]*models.FeedPost, error) {

	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, u.first_name || ' ' || u.last_name AS author_name, u.avatar,
		       p.content, p.media_path, p.media_type, p.privacy, p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE (
			p.user_id IN (
				SELECT following_id FROM followers WHERE follower_id = ?
			)
		)
		AND (
			p.privacy = 'public'
			OR (p.privacy = 'followers' AND EXISTS (
				SELECT 1 FROM followers
				WHERE follower_id = ? AND following_id = p.user_id
			))
			OR (p.privacy = 'selected' AND EXISTS (
				SELECT 1 FROM post_allowed_users
				WHERE post_id = p.id AND user_id = ?
			))
		)
		ORDER BY p.created_at DESC
	`, userID, userID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*models.FeedPost
	for rows.Next() {
		var p models.FeedPost
		err := rows.Scan(&p.ID, &p.UserID, &p.AuthorName, &p.AuthorAvatar, &p.Content, &p.MediaPath, &p.MediaType, &p.Privacy, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &p)
	}
	return posts, nil
}

func (r *postRepository) CreateComment(comment *models.Comment) error {
	_, err := r.db.Exec(`
		INSERT INTO comments (id, post_id, user_id, content, media_path, media_type, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, comment.ID, comment.PostID, comment.UserID, comment.Content, comment.MediaPath, comment.MediaType, comment.CreatedAt)
	return err
}

func (r *postRepository) GetCommentsByPostID(postID string) ([]*models.Comment, error) {

	rows, err := r.db.Query(`
			SELECT id, post_id, user_id, content, media_path, media_type, created_at
			FROM comments
			WHERE post_id = ?
			ORDER BY created_at ASC
		`, postID)

	if err != nil {

		return nil, err
	}

	defer rows.Close()

	var comments []*models.Comment

	for rows.Next() {

		var c models.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.MediaPath, &c.MediaType, &c.CreatedAt)

		if err != nil {

			return nil, err
		}

		comments = append(comments, &c)
	}
	return comments, nil
}