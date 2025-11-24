package stats

type UsersStatsDto struct {
	Count          int `json:"count" db:"count"`
	ActiveCount    int `json:"active_count" db:"active_count"`
	PrAuthorsCount int `json:"pr_authors_count" db:"pr_authors_count"`
}
