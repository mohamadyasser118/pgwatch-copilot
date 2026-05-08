package transform

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Collector reads from the pgwatch sink and produces MetricSignals

type Collector struct {
	pool   *pgxpool.Pool
	sysID  string
	window string // "5m", "1h", "24h"
}

// NewCollector creates a new Collector scoped to one pgwatch instance

func NewCollector(pool *pgxpool.Pool, sysID, window string) *Collector {
	return &Collector{
		pool:   pool,
		sysID:  sysID,
		window: window,
	}
}

// Collect returns MetricSignals for the requested metrics

func (c *Collector) Collect(ctx context.Context, metrics []string) ([]MetricSignal, error) {
	var signals []MetricSignal

	for _, metric := range metrics {
		keys, ok := CounterFields[metric]
		if !ok {
			continue
		}

		for _, key := range keys {
			sig, err := c.computeSignal(ctx, metric, key)
			if err != nil {
				
				continue
			}
			signals = append(signals, sig)
		}
	}

	if len(signals) == 0 {
		return nil, fmt.Errorf(
			"no metric data found for sys_id=%q — check that pgwatch is collecting and the sys_id is correct",
			c.sysID,
		)
	}

	return signals, nil
}

// computeSignal calculates the rate for one specific metric field

func (c *Collector) computeSignal(ctx context.Context, metric, key string) (MetricSignal, error) {

	// Fetch the two most recent samples for this metric + sys_id
	type sample struct {
		T   time.Time
		Val float64
	}

	rows, err := c.pool.Query(ctx, `
		SELECT time, (data->>$1)::float AS val
		FROM metrics.pgwatch2_pgwatch2
		WHERE metric    = $2
		  AND sys_id    = $3
		  AND data->>$1 IS NOT NULL
		  AND (data->>$1)::text ~ '^[0-9]+(\.[0-9]+)?$'
		ORDER BY time DESC
		LIMIT 2
	`, key, metric, c.sysID)
	if err != nil {
		return MetricSignal{}, fmt.Errorf("query %s.%s: %w", metric, key, err)
	}
	defer rows.Close()

	var samples []sample
	for rows.Next() {
		var s sample
		if err := rows.Scan(&s.T, &s.Val); err != nil {
			return MetricSignal{}, fmt.Errorf("scan %s.%s: %w", metric, key, err)
		}
		samples = append(samples, s)
	}

	if len(samples) < 2 {
		return MetricSignal{}, fmt.Errorf("not enough samples for %s.%s (got %d, need 2)", metric, key, len(samples))
	}

	// samples[0] is newest, samples[1] is older
	newest := samples[0]
	older := samples[1]

	elapsed := newest.T.Sub(older.T).Seconds()
	if elapsed <= 0 {
		return MetricSignal{}, fmt.Errorf("invalid time delta for %s.%s: %.2f seconds", metric, key, elapsed)
	}

	delta := newest.Val - older.Val
	if delta < 0 {
		
		return MetricSignal{}, fmt.Errorf("counter reset detected for %s.%s", metric, key)
	}

	currentRate := delta / elapsed

	// Get the baseline rate for this metric (average rate over last hour)
	baselineRate, err := c.getBaselineRate(ctx, metric, key)
	if err != nil || baselineRate <= 0 {
		baselineRate = currentRate
	}

	anomalyScore := 0.0
	if baselineRate > 0 {
		anomalyScore = (currentRate - baselineRate) / baselineRate
	}

	anomalyScore = math.Round(anomalyScore*100) / 100

	return MetricSignal{
		MetricName:   metric + "." + key,
		SysID:        c.sysID,
		CurrentRate:  math.Round(currentRate*10000) / 10000,
		BaselineRate: math.Round(baselineRate*10000) / 10000,
		AnomalyScore: anomalyScore,
		Trend:        TrendFromAnomaly(anomalyScore),
		Window:       c.window,
		CollectedAt:  newest.T,
	}, nil
}

// getBaselineRate calculates the average rate over the last hour
func (c *Collector) getBaselineRate(ctx context.Context, metric, key string) (float64, error) {
	var rate *float64

	err := c.pool.QueryRow(ctx, `
		WITH ordered AS (
			SELECT
				time,
				(data->>$1)::float AS val,
				LAG((data->>$1)::float) OVER (ORDER BY time) AS prev_val,
				EXTRACT(EPOCH FROM (
					time - LAG(time) OVER (ORDER BY time)
				)) AS elapsed_sec
			FROM metrics.pgwatch2_pgwatch2
			WHERE metric  = $2
			  AND sys_id  = $3
			  AND time    > NOW() - INTERVAL '1 hour'
			  AND data->>$1 IS NOT NULL
			  AND (data->>$1)::text ~ '^[0-9]+(\.[0-9]+)?$'
			ORDER BY time
		)
		SELECT AVG((val - prev_val) / NULLIF(elapsed_sec, 0))
		FROM ordered
		WHERE prev_val    IS NOT NULL
		  AND elapsed_sec >  0
		  AND val         >= prev_val  -- skip counter resets
	`, key, metric, c.sysID).Scan(&rate)

	if err != nil {
		return 0, fmt.Errorf("baseline query for %s.%s: %w", metric, key, err)
	}
	if rate == nil {
		return 0, fmt.Errorf("no baseline data for %s.%s", metric, key)
	}

	return *rate, nil
}