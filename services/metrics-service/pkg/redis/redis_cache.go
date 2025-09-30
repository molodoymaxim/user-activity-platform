package redis

import "github.com/molodoymaxim/user-activity-platform/services/metrics-service/internal/domain/entity/metrics"

type metricsCache struct {
	cache Redis
}

func (mr *metricsCache) GetMetricsFromCache(userID metrics.UserID) (*metrics.UserMetrics, error) {
	return nil, nil
}

func (mr *metricsCache) SetMetricsInCache(userMetrics *metrics.UserMetrics) error {
	return nil
}
