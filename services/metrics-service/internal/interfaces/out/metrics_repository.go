package out

import (
	"github.com/molodoymaxim/user-activity-platform/services/metrics-service/internal/domain/entity/metrics"
	"github.com/molodoymaxim/user-activity-platform/services/metrics-service/pkg/postgres"
)

type MetricsRepositoryPort interface {
	GetMetricsFromDB(tx postgres.TxPG, userID metrics.UserID) (*metrics.UserMetrics, error)
	SetMetricsToDB(tx postgres.TxPG, userMetrics *metrics.UserMetrics) error
}
