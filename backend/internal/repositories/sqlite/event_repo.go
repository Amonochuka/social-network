package sqlite

import (
	"database/sql"
	"time"

	"social-network/backend/internal/models"

	"github.com/google/uuid"
)

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *eventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) CreateEvent(e *models.Event) error {
	e.ID = uuid.New().String()
	e.CreatedAt = time.Now()
	_, err := r.db.Exec(
		`INSERT INTO events (id, group_id, creator_id, title, description, event_time, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.GroupID, e.CreatorID, e.Title, e.Description, e.EventTime, e.CreatedAt,
	)
	return err
}

func (r *eventRepository) GetEventsByGroupID(groupID, viewerID string) ([]*models.Event, error) {
	rows, err := r.db.Query(
		`SELECT e.id, e.group_id, e.creator_id, e.title, e.description, e.event_time, e.created_at,
		        (SELECT COUNT(*) FROM event_responses WHERE event_id = e.id AND response = 'going') AS going_count,
		        (SELECT COUNT(*) FROM event_responses WHERE event_id = e.id AND response = 'not_going') AS not_going_count,
		        COALESCE((SELECT response FROM event_responses WHERE event_id = e.id AND user_id = ?), '') AS my_response
		 FROM events e
		 WHERE e.group_id = ?
		 ORDER BY e.event_time ASC`, viewerID, groupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*models.Event
	for rows.Next() {
		e := &models.Event{}
		if err := rows.Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description,
			&e.EventTime, &e.CreatedAt, &e.GoingCount, &e.NotGoing, &e.MyResponse); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (r *eventRepository) GetEventByID(eventID, viewerID string) (*models.Event, error) {
	e := &models.Event{}
	err := r.db.QueryRow(
		`SELECT e.id, e.group_id, e.creator_id, e.title, e.description, e.event_time, e.created_at,
		        (SELECT COUNT(*) FROM event_responses WHERE event_id = e.id AND response = 'going') AS going_count,
		        (SELECT COUNT(*) FROM event_responses WHERE event_id = e.id AND response = 'not_going') AS not_going_count,
		        COALESCE((SELECT response FROM event_responses WHERE event_id = e.id AND user_id = ?), '') AS my_response
		 FROM events e WHERE e.id = ?`, viewerID, eventID,
	).Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description,
		&e.EventTime, &e.CreatedAt, &e.GoingCount, &e.NotGoing, &e.MyResponse)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// UpsertResponse inserts or replaces a user's event response.
// SQLite REPLACE INTO handles the PRIMARY KEY (event_id, user_id) conflict.
func (r *eventRepository) UpsertResponse(resp *models.EventResponse) error {
	resp.CreatedAt = time.Now()
	_, err := r.db.Exec(
		`REPLACE INTO event_responses (event_id, user_id, response, created_at)
		 VALUES (?, ?, ?, ?)`,
		resp.EventID, resp.UserID, resp.Response, resp.CreatedAt,
	)
	return err
}
