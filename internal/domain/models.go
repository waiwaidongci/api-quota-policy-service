package domain

import "time"

type Algorithm string

const (
	FixedWindow   Algorithm = "fixed_window"
	SlidingWindow Algorithm = "sliding_window"
	TokenBucket   Algorithm = "token_bucket"
)

type PolicyStatus string

const (
	Draft      PolicyStatus = "draft"
	Published  PolicyStatus = "published"
	RolledBack PolicyStatus = "rolled_back"
)

type Service struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Environment string            `json:"environment"`
	Tags        map[string]string `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
}

type MatchRule struct {
	ServiceID   string `json:"service_id"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	SourceTag   string `json:"source_tag"`
	Environment string `json:"environment"`
}

type Policy struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Priority      int          `json:"priority"`
	Algorithm     Algorithm    `json:"algorithm"`
	Limit         int64        `json:"limit"`
	WindowSeconds int64        `json:"window_seconds"`
	Burst         int64        `json:"burst"`
	RefillRate    float64      `json:"refill_rate"`
	Rule          MatchRule    `json:"rule"`
	Status        PolicyStatus `json:"status"`
	Version       int          `json:"version"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type DecisionRequest struct {
	ServiceID   string `json:"service_id"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	Source      string `json:"source"`
	Environment string `json:"environment"`
	Key         string `json:"key"`
}
type Decision struct {
	Allowed   bool      `json:"allowed"`
	PolicyID  string    `json:"policy_id"`
	Algorithm string    `json:"algorithm"`
	Key       string    `json:"key"`
	Limit     int64     `json:"limit"`
	Remaining int64     `json:"remaining"`
	ResetAt   time.Time `json:"reset_at"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}
type QuotaSnapshot struct {
	Key       string    `json:"key"`
	PolicyID  string    `json:"policy_id"`
	Count     int64     `json:"count"`
	Limit     int64     `json:"limit"`
	Remaining int64     `json:"remaining"`
	ResetAt   time.Time `json:"reset_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type LimitEvent struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	PolicyID  string    `json:"policy_id"`
	Reason    string    `json:"reason"`
	Count     int64     `json:"count"`
	Limit     int64     `json:"limit"`
	CreatedAt time.Time `json:"created_at"`
}

type PolicyFilter struct {
	Status    *PolicyStatus
	ServiceID string
}
