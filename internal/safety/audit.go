package safety

import (
	"encoding/json"
	"log/slog"
	"time"
)

// AuditEntry records every SQL execution attempt
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	SQL       string    `json:"sql"`
	SysID     string    `json:"sys_id"`
	Action    string    `json:"action"` // Action is "executed" or "rejected".
	Reason    string    `json:"reason,omitempty"`
}

// LogExecution records that a validated SQL was executed.
func LogExecution(sql, sysID string) {
	log(AuditEntry{
		Timestamp: time.Now(),
		SQL:       sql,
		SysID:     sysID,
		Action:    "executed",
	})
}

// LogRejection records that a SQL was rejected by the validator.
func LogRejection(sql, sysID, reason string) {
	log(AuditEntry{
		Timestamp: time.Now(),
		SQL:       sql,
		SysID:     sysID,
		Action:    "rejected",
		Reason:    reason,
	})
}

func log(entry AuditEntry) {
	b, _ := json.Marshal(entry)
	if entry.Action == "rejected" {
		slog.Warn("copilot audit", "entry", string(b))
	} else {
		slog.Info("copilot audit", "entry", string(b))
	}
}