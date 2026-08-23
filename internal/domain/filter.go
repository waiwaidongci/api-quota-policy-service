package domain

import "strings"

func FilterByName(p []Policy, name string) []Policy {
	o := []Policy{}
	name = strings.ToLower(name)
	for _, v := range p {
		if strings.Contains(strings.ToLower(v.Name), name) {
			o = append(o, v)
		}
	}
	return o
}
func FilterByAlgorithm(p []Policy, a Algorithm) []Policy {
	o := []Policy{}
	for _, v := range p {
		if v.Algorithm == a {
			o = append(o, v)
		}
	}
	return o
}
func FilterByPriority(p []Policy, min int) []Policy {
	o := []Policy{}
	for _, v := range p {
		if v.Priority >= min {
			o = append(o, v)
		}
	}
	return o
}
func FilterByEnvironment(s []Service, env string) []Service {
	o := []Service{}
	for _, v := range s {
		if v.Environment == env {
			o = append(o, v)
		}
	}
	return o
}
func MatchAny(p Policy, req DecisionRequest, tags map[string]string) bool {
	return p.Rule.Matches(req, tags)
}
func MatchAll(p []Policy, req DecisionRequest, tags map[string]string) []Policy {
	o := []Policy{}
	for _, v := range p {
		if MatchAny(v, req, tags) {
			o = append(o, v)
		}
	}
	return o
}
