package domain

import "time"

type EventType string

const EventLimitExceeded EventType = "limit_exceeded"

type PolicyEvent struct {
	Type     EventType
	PolicyID string
	Key      string
	At       time.Time
	Details  map[string]string
}
