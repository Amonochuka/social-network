package sqlite

import (
	"database/sql"
	"social-network/backend/internal/models"
	"strings"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, user.ID, user.Email, user.Password, user.FirstName, user.LastName,
		user.DateOfBirth, user.Avatar, user.NickName, user.AboutMe, user.IsPublic)
	return err
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	row := r.db.QueryRow(`
		SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public, created_at
		FROM users WHERE email = ?
	`, email)

	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName,
		&u.DateOfBirth, &u.Avatar, &u.NickName, &u.AboutMe, &u.IsPublic, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetUserByID(id string) (*models.User, error) {
	row := r.db.QueryRow(`
		SELECT id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public, created_at
		FROM users WHERE id = ?
	`, id)

	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName,
		&u.DateOfBirth, &u.Avatar, &u.NickName, &u.AboutMe, &u.IsPublic, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) SearchUsers(query, currentUserID string) ([]*models.UserSearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []*models.UserSearchResult{}, nil
	}

	search := "%" + strings.ToLower(query) + "%"

	rows, err := r.db.Query(`
		SELECT
			id,
			first_name,
			last_name,
			COALESCE(nickname, ''),
			COALESCE(avatar, ''),
			NOT is_public
		FROM users
		WHERE id != ?
		AND (
			LOWER(first_name) LIKE ?
			OR LOWER(last_name) LIKE ?
			OR LOWER(COALESCE(nickname, '')) LIKE ?
			OR LOWER(email) LIKE ?
		)
		ORDER BY
			CASE
				WHEN LOWER(first_name) = LOWER(?) THEN 1
				WHEN LOWER(last_name) = LOWER(?) THEN 2
				WHEN LOWER(COALESCE(nickname, '')) = LOWER(?) THEN 3
				ELSE 4
			END,
			first_name,
			last_name
		LIMIT 20
	`,
		currentUserID,
		search,
		search,
		search,
		search,
		query,
		query,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.UserSearchResult{}

	for rows.Next() {
		u := &models.UserSearchResult{}

		if err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Nickname,
			&u.Avatar,
			&u.IsPrivate,
		); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET
			email = ?,
			first_name = ?,
			last_name = ?,
			date_of_birth = ?,
			nickname = ?,
			about_me = ?
		WHERE id = ?
	`,
		user.Email,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		user.NickName,
		user.AboutMe,
		user.ID,
	)
	return err
}

func (r *userRepository) UpdateAvatar(userID, avatarPath string) error {
	_, err := r.db.Exec(`
		UPDATE users SET avatar = ? WHERE id = ?
	`, avatarPath, userID)
	return err
}

func (r *userRepository) UpdatePrivacy(userID string, isPublic bool) error {
	_, err := r.db.Exec(`
		UPDATE users SET is_public = ? WHERE id = ?
	`, isPublic, userID)
	return err
}

func (r *userRepository) UpdatePassword(userID, hashedPassword string) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET password = ?
		WHERE id = ?
	`, hashedPassword, userID)

	return err
}