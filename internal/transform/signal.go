package transform

import "time"

// MetricSignal is the output of the transformation layer.

type MetricSignal struct {
	MetricName   string   
	SysID        string    
	CurrentRate  float64   
	BaselineRate float64   
	AnomalyScore float64   
	Trend        string    
	Window       string    
	CollectedAt  time.Time 
}


func TrendFromAnomaly(score float64) string {
	switch {
	case score >= 0.5:
		return "spike"
	case score < -0.3:
		return "declining"
	default:
		return "stable"
	}
}