package models

type UserSearchResult struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsPrivate bool   `json:"is_private"`
}

type GroupSearchResult struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SearchResponse struct {
	Users  []*UserSearchResult  `json:"users"`
	Groups []*GroupSearchResult `json:"groups"`
}