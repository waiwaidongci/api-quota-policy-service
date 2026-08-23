package domain

import "strings"

func (r MatchRule) Matches(req DecisionRequest, tags map[string]string) bool {
	if r.ServiceID != "" && r.ServiceID != req.ServiceID {
		return false
	}
	if r.Environment != "" && r.Environment != req.Environment {
		return false
	}
	if r.Method != "" && !strings.EqualFold(r.Method, req.Method) {
		return false
	}
	if r.Path != "" && !pathMatches(r.Path, req.Path) {
		return false
	}
	if r.SourceTag != "" && tags[req.Source] != r.SourceTag {
		return false
	}
	return true
}

func pathMatches(pattern, path string) bool {
	if pattern == path || pattern == "" {
		return true
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(path, strings.TrimSuffix(pattern, "*"))
	}
	return false
}
