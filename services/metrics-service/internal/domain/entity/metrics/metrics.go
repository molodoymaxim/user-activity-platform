package metrics

type UserID int

type UserMetrics struct {
	UserID    UserID `json:"user_id"`
	Clicks    int    `json:"clicks"`
	PageViews int    `json:"page_views"`
	LastSeen  int    `json:"last_seen"`
}
