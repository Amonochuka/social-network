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

func (r *postRepository) GetPostsByUserID(userID, viewerID string) ([]*models.FeedPost, error) {
	rows, err := r.db.Query(`
		SELECT
			p.id,
			p.user_id,
			u.first_name || ' ' || u.last_name AS author_name,
			u.avatar,
			p.content,
			p.media_path,
			p.media_type,
			p.privacy,
			COALESCE((
				SELECT COUNT(*)
				FROM comments c
				WHERE c.post_id = p.id
			), 0) AS comment_count,
			COALESCE((
				SELECT COUNT(*)
				FROM post_likes l
				WHERE l.post_id = p.id
			), 0) AS like_count,
			EXISTS(
				SELECT 1
				FROM post_likes l
				WHERE l.post_id = p.id
				AND l.user_id = ?
			) AS liked_by_me,
			p.created_at,
			p.updated_at
		FROM posts p
		JOIN users u
			ON u.id = p.user_id
		WHERE p.user_id = ?
		AND (
			p.user_id = ?
			OR p.privacy = 'public'
			OR (
				p.privacy = 'followers'
				AND EXISTS (
					SELECT 1
					FROM followers
					WHERE follower_id = ?
					AND following_id = ?
				)
			)
			OR (
				p.privacy = 'selected'
				AND EXISTS (
					SELECT 1
					FROM post_allowed_users
					WHERE post_id = p.id
					AND user_id = ?
				)
			)
		)
		ORDER BY p.created_at DESC
	`,
		viewerID,
		userID,
		viewerID,
		viewerID,
		userID,
		viewerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.FeedPost

	for rows.Next() {
		var p models.FeedPost

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.AuthorName,
			&p.AuthorAvatar,
			&p.Content,
			&p.MediaPath,
			&p.MediaType,
			&p.Privacy,
			&p.CommentCount,
			&p.LikeCount,
			&p.LikedByMe,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &p)
	}

	return posts, nil
}

// GetFeed joins users to include author name, avatar, like count, and liked-by-current-user flag.
func (r *postRepository) GetFeed(userID string) ([]*models.FeedPost, error) {

	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, u.first_name || ' ' || u.last_name AS author_name, u.avatar,
		       p.content, p.media_path, p.media_type, p.privacy,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count,
		       (SELECT COUNT(*) FROM post_likes pl WHERE pl.post_id = p.id) AS like_count,
		       (SELECT COUNT(*) FROM post_likes pl WHERE pl.post_id = p.id AND pl.user_id = ?) AS liked_by_me,
		       p.created_at, p.updated_at
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE (
			p.user_id = ?
			OR p.user_id IN (
				SELECT following_id FROM followers WHERE follower_id = ?
			)
		)
		AND (
			p.user_id = ?
			OR p.privacy = 'public'
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
	`, userID, userID, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*models.FeedPost
	for rows.Next() {
		var p models.FeedPost
		var likedByMe int
		err := rows.Scan(
			&p.ID, &p.UserID, &p.AuthorName, &p.AuthorAvatar,
			&p.Content, &p.MediaPath, &p.MediaType, &p.Privacy,
			&p.CommentCount, &p.LikeCount, &likedByMe,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		p.LikedByMe = likedByMe > 0
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

func (r *postRepository) GetCommentsByPostID(postID string) ([]*models.CommentDetail, error) {

	rows, err := r.db.Query(`
			SELECT c.id, c.post_id, c.user_id, u.first_name || ' ' || u.last_name AS author_name, u.avatar,
			       c.content, c.media_path, c.media_type, c.created_at
			FROM comments c
			JOIN users u ON u.id = c.user_id
			WHERE c.post_id = ?
			ORDER BY c.created_at ASC
		`, postID)

	if err != nil {

		return nil, err
	}

	defer rows.Close()

	var comments []*models.CommentDetail

	for rows.Next() {

		var c models.CommentDetail
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.AuthorName, &c.AuthorAvatar, &c.Content, &c.MediaPath, &c.MediaType, &c.CreatedAt)

		if err != nil {

			return nil, err
		}

		comments = append(comments, &c)
	}
	return comments, nil
}

// ToggleLike adds a like if the user hasn't liked yet, or removes it if they have.
// Returns the new liked state and the updated total like count.
func (r *postRepository) ToggleLike(postID, userID string) (bool, int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND user_id = ?`,
		postID, userID,
	).Scan(&count)
	if err != nil {
		return false, 0, err
	}

	var liked bool
	if count > 0 {
		// already liked — remove it
		_, err = r.db.Exec(`DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`, postID, userID)
		liked = false
	} else {
		// not yet liked — add it
		_, err = r.db.Exec(
			`INSERT INTO post_likes (post_id, user_id, created_at) VALUES (?, ?, datetime('now'))`,
			postID, userID,
		)
		liked = true
	}
	if err != nil {
		return false, 0, err
	}

	var total int
	err = r.db.QueryRow(`SELECT COUNT(*) FROM post_likes WHERE post_id = ?`, postID).Scan(&total)
	if err != nil {
		return liked, 0, err
	}
	return liked, total, nil
}

func (r *postRepository) HasLiked(postID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND user_id = ?`,
		postID, userID,
	).Scan(&count)
	return count > 0, err
}

func (r *postRepository) CountLikes(postID string) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM post_likes WHERE post_id = ?`, postID).Scan(&count)
	return count, err
}
