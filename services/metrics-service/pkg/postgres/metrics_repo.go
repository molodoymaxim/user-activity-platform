package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/molodoymaxim/user-activity-platform/services/metrics-service/internal/domain/entity/metrics"
	"time"
)

type metricsRepo struct {
	db Postgre
}

func (mr *metricsRepo) GetMetricsFromDB(tx TxPG, userID metrics.UserID) (*metrics.UserMetrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	query := `
		SELECT user_id, clicks, page_views, last_seen
		FROM user_metrics
		WHERE user_id = $1;
	`

	var (
		err         error
		userMetrics metrics.UserMetrics
	)

	if tx == nil {
		err = mr.db.QueryRow(ctx, query, userID).Scan(
			&userMetrics.UserID,
			&userMetrics.Clicks,
			&userMetrics.PageViews,
			&userMetrics.LastSeen,
		)
	}

	if tx != nil {
		err = tx.QueryRow(ctx, query, userID).Scan(
			&userMetrics.UserID,
			&userMetrics.Clicks,
			&userMetrics.PageViews,
			&userMetrics.LastSeen,
		)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get metrics from db: %w", err)
	}

	return &userMetrics, nil
}

func (mr *metricsRepo) SetMetricsToDB(tx TxPG, userMetrics *metrics.UserMetrics) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	query := `
		INSERT INTO user_metrics(clicks, page_views, last_seen)
		VALUES ($1, $2, $3);
	`

	var (
		err error
	)

	if tx == nil {
		_, err = mr.db.Exec(ctx, query, userMetrics.Clicks, userMetrics.PageViews, userMetrics.LastSeen)
	}

	if tx != nil {
		_, err = tx.Exec(ctx, query, userMetrics.Clicks, userMetrics.PageViews, userMetrics.LastSeen)
	}

	if err != nil {
		return fmt.Errorf("failed to set metrics to db: %w", err)
	}

	return nil
}
