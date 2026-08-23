package domain

import (
	"encoding/json"
	"fmt"
)

func EncodePolicy(p Policy) ([]byte, error) { return json.Marshal(p) }
func DecodePolicy(b []byte) (Policy, error) {
	var p Policy
	if err := json.Unmarshal(b, &p); err != nil {
		return Policy{}, fmt.Errorf("decode policy: %w", err)
	}
	return p, nil
}
func EncodeDecision(d Decision) ([]byte, error) { return json.Marshal(d) }
func DecodeDecision(b []byte) (Decision, error) {
	var d Decision
	if err := json.Unmarshal(b, &d); err != nil {
		return Decision{}, fmt.Errorf("decode decision: %w", err)
	}
	return d, nil
}
func EncodeEvent(e LimitEvent) ([]byte, error) { return json.Marshal(e) }
func DecodeEvent(b []byte) (LimitEvent, error) {
	var e LimitEvent
	if err := json.Unmarshal(b, &e); err != nil {
		return LimitEvent{}, fmt.Errorf("decode event: %w", err)
	}
	return e, nil
}
func ClonePolicy(p Policy) Policy { q := p; q.Rule = p.Rule; return q }
func CloneService(s Service) Service {
	q := s
	q.Tags = map[string]string{}
	for k, v := range s.Tags {
		q.Tags[k] = v
	}
	return q
}
