package sqlite

import (
	"database/sql"
	"time"

	"social-network/backend/internal/models"

	"github.com/google/uuid"
)

type groupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *groupRepository {
	return &groupRepository{db: db}
}

// ── Core group CRUD ──────────────────────────────────────────────────────────

func (r *groupRepository) CreateGroup(g *models.Group) error {
	g.ID = uuid.New().String()
	g.CreatedAt = time.Now()
	_, err := r.db.Exec(
		`INSERT INTO groups (id, creator_id, title, description, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		g.ID, g.CreatorID, g.Title, g.Description, g.CreatedAt,
	)
	return err
}

func (r *groupRepository) GetGroupByID(groupID string) (*models.Group, error) {
	g := &models.Group{}
	err := r.db.QueryRow(
		`SELECT id, creator_id, title, description, created_at,
		        (SELECT COUNT(*) FROM group_members WHERE group_id = groups.id) AS member_count
		 FROM groups WHERE id = ?`, groupID,
	).Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &g.CreatedAt, &g.MemberCount)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *groupRepository) GetAllGroups(viewerID string) ([]*models.Group, error) {
	rows, err := r.db.Query(
		`SELECT g.id, g.creator_id, g.title, g.description, g.created_at,
		        (SELECT COUNT(*) FROM group_members WHERE group_id = g.id) AS member_count,
		        EXISTS(SELECT 1 FROM group_members WHERE group_id = g.id AND user_id = ?) AS is_member
		 FROM groups g
		 ORDER BY g.created_at DESC`, viewerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []*models.Group
	for rows.Next() {
		g := &models.Group{}
		if err := rows.Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &g.CreatedAt, &g.MemberCount, &g.IsMember); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

// ── Membership ───────────────────────────────────────────────────────────────

func (r *groupRepository) AddMember(groupID, userID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO group_members (group_id, user_id, joined_at) VALUES (?, ?, ?)`,
		groupID, userID, time.Now(),
	)
	return err
}

func (r *groupRepository) RemoveMember(groupID, userID string) error {
	_, err := r.db.Exec(
		`DELETE FROM group_members WHERE group_id = ? AND user_id = ?`,
		groupID, userID,
	)
	return err
}

func (r *groupRepository) IsMember(groupID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?`,
		groupID, userID,
	).Scan(&count)
	return count > 0, err
}

func (r *groupRepository) GetMembers(groupID string) ([]*models.GroupMember, error) {
	rows, err := r.db.Query(
		`SELECT group_id, user_id, joined_at FROM group_members WHERE group_id = ? ORDER BY joined_at ASC`,
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []*models.GroupMember
	for rows.Next() {
		m := &models.GroupMember{}
		if err := rows.Scan(&m.GroupID, &m.UserID, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

// ── Invitations ──────────────────────────────────────────────────────────────

func (r *groupRepository) CreateInvitation(inv *models.GroupInvitation) error {
	inv.ID = uuid.New().String()
	inv.Status = "pending"
	inv.CreatedAt = time.Now()
	_, err := r.db.Exec(
		`INSERT INTO group_invitations (id, group_id, inviter_id, invitee_id, status, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		inv.ID, inv.GroupID, inv.InviterID, inv.InviteeID, inv.Status, inv.CreatedAt,
	)
	return err
}

func (r *groupRepository) GetInvitationByID(invID string) (*models.GroupInvitation, error) {
	inv := &models.GroupInvitation{}
	err := r.db.QueryRow(
		`SELECT id, group_id, inviter_id, invitee_id, status, created_at
		 FROM group_invitations WHERE id = ?`, invID,
	).Scan(&inv.ID, &inv.GroupID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *groupRepository) GetPendingInvitation(groupID, inviteeID string) (*models.GroupInvitation, error) {
	inv := &models.GroupInvitation{}
	err := r.db.QueryRow(
		`SELECT id, group_id, inviter_id, invitee_id, status, created_at
		 FROM group_invitations WHERE group_id = ? AND invitee_id = ? AND status = 'pending'`,
		groupID, inviteeID,
	).Scan(&inv.ID, &inv.GroupID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *groupRepository) UpdateInvitation(invID, status string) error {
	_, err := r.db.Exec(
		`UPDATE group_invitations SET status = ? WHERE id = ?`, status, invID,
	)
	return err
}

// ── Join requests ────────────────────────────────────────────────────────────

func (r *groupRepository) CreateJoinRequest(req *models.GroupRequest) error {
	req.ID = uuid.New().String()
	req.Status = "pending"
	req.CreatedAt = time.Now()
	_, err := r.db.Exec(
		`INSERT INTO group_requests (id, group_id, user_id, status, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		req.ID, req.GroupID, req.UserID, req.Status, req.CreatedAt,
	)
	return err
}

func (r *groupRepository) GetJoinRequestByID(reqID string) (*models.GroupRequest, error) {
	req := &models.GroupRequest{}
	err := r.db.QueryRow(
		`SELECT id, group_id, user_id, status, created_at FROM group_requests WHERE id = ?`, reqID,
	).Scan(&req.ID, &req.GroupID, &req.UserID, &req.Status, &req.CreatedAt)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *groupRepository) GetPendingJoinRequest(groupID, userID string) (*models.GroupRequest, error) {
	req := &models.GroupRequest{}
	err := r.db.QueryRow(
		`SELECT id, group_id, user_id, status, created_at
		 FROM group_requests WHERE group_id = ? AND user_id = ? AND status = 'pending'`,
		groupID, userID,
	).Scan(&req.ID, &req.GroupID, &req.UserID, &req.Status, &req.CreatedAt)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *groupRepository) UpdateJoinRequest(reqID, status string) error {
	_, err := r.db.Exec(
		`UPDATE group_requests SET status = ? WHERE id = ?`, status, reqID,
	)
	return err
}

// ── Group posts ──────────────────────────────────────────────────────────────

func (r *groupRepository) CreateGroupPost(post *models.GroupPost) error {
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()
	_, err := r.db.Exec(
		`INSERT INTO group_posts (id, group_id, user_id, content, media_path, media_type, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		post.ID, post.GroupID, post.UserID, post.Content, post.MediaPath, post.MediaType, post.CreatedAt,
	)
	return err
}

func (r *groupRepository) GetGroupPosts(groupID string) ([]*models.GroupPost, error) {
	rows, err := r.db.Query(
		`SELECT id, group_id, user_id, content, media_path, media_type, created_at
		 FROM group_posts WHERE group_id = ? ORDER BY created_at DESC`, groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*models.GroupPost
	for rows.Next() {
		p := &models.GroupPost{}
		if err := rows.Scan(&p.ID, &p.GroupID, &p.UserID, &p.Content, &p.MediaPath, &p.MediaType, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *groupRepository) GetGroupPostByID(postID string) (*models.GroupPost, error) {
	p := &models.GroupPost{}
	err := r.db.QueryRow(
		`SELECT id, group_id, user_id, content, media_path, media_type, created_at
		 FROM group_posts WHERE id = ?`, postID,
	).Scan(&p.ID, &p.GroupID, &p.UserID, &p.Content, &p.MediaPath, &p.MediaType, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ── Group comments ───────────────────────────────────────────────────────────

func (r *groupRepository) CreateGroupComment(comment *models.GroupComment) error {
	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now()
	_, err := r.db.Exec(
		`INSERT INTO group_comments (id, group_post_id, user_id, content, media_path, media_type, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		comment.ID, comment.GroupPostID, comment.UserID, comment.Content,
		comment.MediaPath, comment.MediaType, comment.CreatedAt,
	)
	return err
}

func (r *groupRepository) GetGroupComments(groupPostID string) ([]*models.GroupComment, error) {
	rows, err := r.db.Query(
		`SELECT id, group_post_id, user_id, content, media_path, media_type, created_at
		 FROM group_comments WHERE group_post_id = ? ORDER BY created_at ASC`, groupPostID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.GroupComment
	for rows.Next() {
		c := &models.GroupComment{}
		if err := rows.Scan(&c.ID, &c.GroupPostID, &c.UserID, &c.Content, &c.MediaPath, &c.MediaType, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
