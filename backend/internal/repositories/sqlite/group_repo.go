package sqlite

import (
	"database/sql"
	"time"

	"social-network/backend/internal/models"

	"github.com/google/uuid"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

// ── Groups ──────────────────────────────────────────────────────────────────

func (r *GroupRepository) CreateGroup(group *models.Group) error {
	_, err := r.db.Exec(
		`INSERT INTO groups (id, creator_id, title, description, created_at) VALUES (?, ?, ?, ?, ?)`,
		group.ID, group.CreatorID, group.Title, group.Description, group.CreatedAt,
	)
	return err
}

func (r *GroupRepository) GetGroupByID(groupID string) (*models.Group, error) {
	row := r.db.QueryRow(`SELECT id, creator_id, title, description, created_at FROM groups WHERE id = ?`, groupID)
	g := &models.Group{}
	return g, row.Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &g.CreatedAt)
}

func (r *GroupRepository) GetAllGroups() ([]*models.Group, error) {
	rows, err := r.db.Query(`SELECT id, creator_id, title, description, created_at FROM groups ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []*models.Group
	for rows.Next() {
		g := &models.Group{}
		if err := rows.Scan(&g.ID, &g.CreatorID, &g.Title, &g.Description, &g.CreatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// ── Members ─────────────────────────────────────────────────────────────────

func (r *GroupRepository) AddMember(groupID, userID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO group_members (group_id, user_id, joined_at) VALUES (?, ?, ?)`,
		groupID, userID, time.Now(),
	)
	return err
}

func (r *GroupRepository) RemoveMember(groupID, userID string) error {
	_, err := r.db.Exec(`DELETE FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, userID)
	return err
}

func (r *GroupRepository) IsMember(groupID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, userID).Scan(&count)
	return count > 0, err
}

func (r *GroupRepository) GetMembers(groupID string) ([]*models.MemberProfile, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.first_name, u.last_name, COALESCE(u.avatar,''), COALESCE(u.nickname,'')
		FROM group_members gm
		JOIN users u ON u.id = gm.user_id
		WHERE gm.group_id = ?
		ORDER BY gm.joined_at`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []*models.MemberProfile
	for rows.Next() {
		m := &models.MemberProfile{}
		if err := rows.Scan(&m.UserID, &m.FirstName, &m.LastName, &m.Avatar, &m.Nickname); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *GroupRepository) GetMemberIDs(groupID string) ([]string, error) {
	rows, err := r.db.Query(`SELECT user_id FROM group_members WHERE group_id = ?`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// ── Invitations ──────────────────────────────────────────────────────────────

func (r *GroupRepository) CreateInvitation(inv *models.GroupInvitation) error {
	_, err := r.db.Exec(
		`INSERT INTO group_invitations (id, group_id, inviter_id, invitee_id, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		inv.ID, inv.GroupID, inv.InviterID, inv.InviteeID, inv.Status, inv.CreatedAt,
	)
	return err
}

func (r *GroupRepository) GetInvitationByID(invID string) (*models.GroupInvitation, error) {
	row := r.db.QueryRow(`SELECT id, group_id, inviter_id, invitee_id, status, created_at FROM group_invitations WHERE id = ?`, invID)
	inv := &models.GroupInvitation{}
	return inv, row.Scan(&inv.ID, &inv.GroupID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt)
}

func (r *GroupRepository) GetPendingInvitation(groupID, inviteeID string) (*models.GroupInvitation, error) {
	row := r.db.QueryRow(
		`SELECT id, group_id, inviter_id, invitee_id, status, created_at FROM group_invitations WHERE group_id = ? AND invitee_id = ? AND status = 'pending'`,
		groupID, inviteeID,
	)
	inv := &models.GroupInvitation{}
	err := row.Scan(&inv.ID, &inv.GroupID, &inv.InviterID, &inv.InviteeID, &inv.Status, &inv.CreatedAt)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *GroupRepository) UpdateInvitationStatus(invID, status string) error {
	_, err := r.db.Exec(`UPDATE group_invitations SET status = ? WHERE id = ?`, status, invID)
	return err
}

func (r *GroupRepository) GetInvitationStatus(groupID, userID string) (string, error) {
	var status string
	err := r.db.QueryRow(
		`SELECT status FROM group_invitations WHERE group_id = ? AND invitee_id = ? ORDER BY created_at DESC LIMIT 1`,
		groupID, userID,
	).Scan(&status)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return status, err
}

// ── Join Requests ────────────────────────────────────────────────────────────

func (r *GroupRepository) CreateJoinRequest(req *models.GroupRequest) error {
	_, err := r.db.Exec(
		`INSERT INTO group_requests (id, group_id, user_id, status, created_at) VALUES (?, ?, ?, ?, ?)`,
		req.ID, req.GroupID, req.UserID, req.Status, req.CreatedAt,
	)
	return err
}

func (r *GroupRepository) GetJoinRequestByID(reqID string) (*models.GroupRequest, error) {
	row := r.db.QueryRow(`SELECT id, group_id, user_id, status, created_at FROM group_requests WHERE id = ?`, reqID)
	req := &models.GroupRequest{}
	return req, row.Scan(&req.ID, &req.GroupID, &req.UserID, &req.Status, &req.CreatedAt)
}

func (r *GroupRepository) GetPendingJoinRequest(groupID, userID string) (*models.GroupRequest, error) {
	row := r.db.QueryRow(
		`SELECT id, group_id, user_id, status, created_at FROM group_requests WHERE group_id = ? AND user_id = ? AND status = 'pending'`,
		groupID, userID,
	)
	req := &models.GroupRequest{}
	err := row.Scan(&req.ID, &req.GroupID, &req.UserID, &req.Status, &req.CreatedAt)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (r *GroupRepository) UpdateJoinRequestStatus(reqID, status string) error {
	_, err := r.db.Exec(`UPDATE group_requests SET status = ? WHERE id = ?`, status, reqID)
	return err
}

func (r *GroupRepository) GetJoinRequestStatus(groupID, userID string) (string, error) {
	var status string
	err := r.db.QueryRow(
		`SELECT status FROM group_requests WHERE group_id = ? AND user_id = ? ORDER BY created_at DESC LIMIT 1`,
		groupID, userID,
	).Scan(&status)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return status, err
}

// ── Events ───────────────────────────────────────────────────────────────────

func (r *GroupRepository) CreateEvent(event *models.Event) error {
	_, err := r.db.Exec(
		`INSERT INTO events (id, group_id, creator_id, title, description, event_time, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.GroupID, event.CreatorID, event.Title, event.Description, event.EventTime, event.CreatedAt,
	)
	return err
}

func (r *GroupRepository) GetEventsByGroupID(groupID string) ([]*models.Event, error) {
	rows, err := r.db.Query(
		`SELECT id, group_id, creator_id, title, description, event_time, created_at FROM events WHERE group_id = ? ORDER BY event_time`,
		groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*models.Event
	for rows.Next() {
		e := &models.Event{}
		if err := rows.Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventTime, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *GroupRepository) GetEventByID(eventID string) (*models.Event, error) {
	row := r.db.QueryRow(`SELECT id, group_id, creator_id, title, description, event_time, created_at FROM events WHERE id = ?`, eventID)
	e := &models.Event{}
	return e, row.Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventTime, &e.CreatedAt)
}

func (r *GroupRepository) UpsertEventResponse(eventID, userID, response string) error {
	_, err := r.db.Exec(
		`INSERT INTO event_responses (event_id, user_id, response, created_at) VALUES (?, ?, ?, ?)
		 ON CONFLICT(event_id, user_id) DO UPDATE SET response = excluded.response`,
		eventID, userID, response, time.Now(),
	)
	return err
}

func (r *GroupRepository) GetEventResponses(eventID string) (going int, notGoing int, err error) {
	rows, err := r.db.Query(`SELECT response FROM event_responses WHERE event_id = ?`, eventID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var resp string
		if err := rows.Scan(&resp); err != nil {
			return 0, 0, err
		}
		if resp == "going" {
			going++
		} else {
			notGoing++
		}
	}
	return going, notGoing, nil
}

func (r *GroupRepository) GetMyEventResponse(eventID, userID string) (string, error) {
	var response string
	err := r.db.QueryRow(`SELECT response FROM event_responses WHERE event_id = ? AND user_id = ?`, eventID, userID).Scan(&response)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return response, err
}

// ── Group Posts ──────────────────────────────────────────────────────────────

func (r *GroupRepository) CreateGroupPost(post *models.GroupPost) error {
	_, err := r.db.Exec(
		`INSERT INTO group_posts (id, group_id, user_id, content, media_path, media_type, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		post.ID, post.GroupID, post.UserID, post.Content, post.MediaPath, post.MediaType, post.CreatedAt,
	)
	return err
}

func (r *GroupRepository) GetGroupPosts(groupID string) ([]*models.GroupPostDetail, error) {
	rows, err := r.db.Query(`
		SELECT gp.id, gp.group_id, gp.user_id,
		       u.first_name || ' ' || u.last_name,
		       COALESCE(u.avatar,''),
		       gp.content, COALESCE(gp.media_path,''), COALESCE(gp.media_type,''),
		       (SELECT COUNT(*) FROM group_comments gc WHERE gc.group_post_id = gp.id),
		       gp.created_at
		FROM group_posts gp
		JOIN users u ON u.id = gp.user_id
		WHERE gp.group_id = ?
		ORDER BY gp.created_at DESC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*models.GroupPostDetail
	for rows.Next() {
		p := &models.GroupPostDetail{}
		if err := rows.Scan(&p.ID, &p.GroupID, &p.UserID, &p.AuthorName, &p.AuthorAvatar,
			&p.Content, &p.MediaPath, &p.MediaType, &p.CommentCount, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *GroupRepository) CreateGroupComment(comment *models.GroupComment) error {
	_, err := r.db.Exec(
		`INSERT INTO group_comments (id, group_post_id, user_id, content, media_path, media_type, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		comment.ID, comment.GroupPostID, comment.UserID, comment.Content, comment.MediaPath, comment.MediaType, comment.CreatedAt,
	)
	return err
}

func (r *GroupRepository) GetGroupComments(groupPostID string) ([]*models.GroupCommentDetail, error) {
	rows, err := r.db.Query(`
		SELECT gc.id, gc.group_post_id, gc.user_id,
		       u.first_name || ' ' || u.last_name,
		       COALESCE(u.avatar,''),
		       gc.content, COALESCE(gc.media_path,''), COALESCE(gc.media_type,''), gc.created_at
		FROM group_comments gc
		JOIN users u ON u.id = gc.user_id
		WHERE gc.group_post_id = ?
		ORDER BY gc.created_at`, groupPostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.GroupCommentDetail
	for rows.Next() {
		c := &models.GroupCommentDetail{}
		if err := rows.Scan(&c.ID, &c.GroupPostID, &c.UserID, &c.AuthorName, &c.AuthorAvatar,
			&c.Content, &c.MediaPath, &c.MediaType, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// ── Group Chat ───────────────────────────────────────────────────────────────

func (r *GroupRepository) CreateGroupChatMessage(msg *models.GroupMessage) error {
	_, err := r.db.Exec(
		`INSERT INTO group_messages (id, group_id, sender_id, content, created_at) VALUES (?, ?, ?, ?, ?)`,
		msg.ID, msg.GroupID, msg.SenderID, msg.Content, msg.CreatedAt,
	)
	return err
}

func (r *GroupRepository) GetGroupChatMessages(groupID string) ([]*models.GroupMessageDetailFull, error) {
	rows, err := r.db.Query(`
		SELECT gm.id, gm.group_id, gm.sender_id,
		       u.first_name || ' ' || u.last_name,
		       COALESCE(u.avatar,''),
		       gm.content, gm.created_at
		FROM group_messages gm
		JOIN users u ON u.id = gm.sender_id
		WHERE gm.group_id = ?
		ORDER BY gm.created_at`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []*models.GroupMessageDetailFull
	for rows.Next() {
		m := &models.GroupMessageDetailFull{}
		if err := rows.Scan(&m.ID, &m.GroupID, &m.SenderID, &m.SenderName, &m.SenderAvatar,
			&m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// ── Unused interface stubs (satisfy ChatRepository interface) ────────────────

func (r *GroupRepository) CreateGroupMessage(msg *models.GroupMessage) error {
	return r.CreateGroupChatMessage(msg)
}

func (r *GroupRepository) GetGroupMessages(groupID string) ([]*models.GroupMessageDetail, error) {
	// Returns the basic type for backward compatibility — use GetGroupChatMessages for full detail.
	rows, err := r.db.Query(`SELECT id, group_id, sender_id, content, created_at FROM group_messages WHERE group_id = ? ORDER BY created_at`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var msgs []*models.GroupMessageDetail
	for rows.Next() {
		m := &models.GroupMessageDetail{}
		if err := rows.Scan(&m.ID, &m.GroupID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// NewGroupID returns a new UUID string.
func NewGroupID() string {
	return uuid.New().String()
}
