package out

import "github.com/molodoymaxim/user-activity-platform/services/metrics-service/internal/domain/entity/metrics"

type MetricsCachePort interface {
	GetMetricsFromCache(userID metrics.UserID) (*metrics.UserMetrics, error)
	SetMetricsInCache(userMetrics *metrics.UserMetrics) error
}
