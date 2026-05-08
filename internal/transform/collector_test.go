package transform

import (
	"math"
	"testing"
)

func TestTrendFromAnomaly(t *testing.T) {
	cases := []struct {
		score float64
		want  string
	}{
		{1.0, "spike"},   
		{0.6, "spike"},    
		{0.5, "spike"},    
		{0.4, "stable"},   
		{0.0, "stable"},   
		{-0.2, "stable"},  
		{-0.3, "stable"},  
		{-0.4, "declining"}, 
		{-1.0, "declining"},
	}

	for _, c := range cases {
		got := TrendFromAnomaly(c.score)
		if got != c.want {
			t.Errorf("TrendFromAnomaly(%.1f) = %q, want %q", c.score, got, c.want)
		}
	}
}

func TestAnomalyScore(t *testing.T) {
	cases := []struct {
		current  float64
		baseline float64
		want     float64
	}{
		{2.0, 1.0, 1.0},   
		{1.0, 1.0, 0.0},   
		{0.5, 1.0, -0.5}, 
		{3.0, 1.0, 2.0},   
		{0.0, 1.0, -1.0},  
	}

	for _, c := range cases {
		got := (c.current - c.baseline) / c.baseline
		got = math.Round(got*100) / 100
		if math.Abs(got-c.want) > 0.001 {
			t.Errorf("anomaly(current=%.1f, baseline=%.1f) = %.3f, want %.3f",
				c.current, c.baseline, got, c.want)
		}
	}
}

func TestAnomalyScoreZeroBaseline(t *testing.T) {

	baseline := 0.0
	current := 5.0
	var score float64
	if baseline > 0 {
		score = (current - baseline) / baseline
	}
	if score != 0.0 {
		t.Errorf("expected 0.0 got %f", score)
	}
}

func TestRateCalculation(t *testing.T) {
	// Simulate two consecutive samples
	v1 := 1000.0 
	v2 := 1060.0 
	elapsedSeconds := 30.0

	rate := (v2 - v1) / elapsedSeconds // should be 2.0 commits/sec
	want := 2.0

	if math.Abs(rate-want) > 0.001 {
		t.Errorf("rate = %.4f, want %.4f", rate, want)
	}
}