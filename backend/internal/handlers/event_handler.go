package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"social-network/backend/internal/middleware"
	"social-network/backend/internal/services"
)

type EventHandler struct {
	eventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

// POST /api/groups/{group_id}/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	groupID := r.PathValue("group_id")
	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		EventTime   string `json:"event_time"` // RFC3339 e.g. "2025-08-01T14:00:00Z"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	eventTime, err := time.Parse(time.RFC3339, body.EventTime)
	if err != nil {
		http.Error(w, "event_time must be RFC3339 format e.g. 2025-08-01T14:00:00Z", http.StatusBadRequest)
		return
	}
	event, err := h.eventService.CreateEvent(groupID, userID, body.Title, body.Description, eventTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// GET /api/groups/{group_id}/events
func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	groupID := r.PathValue("group_id")
	events, err := h.eventService.GetEventsByGroup(groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// POST /api/events/{event_id}/respond
func (h *EventHandler) RespondToEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := middleware.GetUserID(r)
	eventID := r.PathValue("event_id")
	var body struct {
		Response string `json:"response"` // "going" | "not_going"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.eventService.RespondToEvent(eventID, userID, body.Response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "response recorded"})
}
