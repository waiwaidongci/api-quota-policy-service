package application

import (
	"fmt"
	"github.com/example/api-quota-service/internal/domain"
	"regexp"
	"strings"
)

type RuleCompiler struct{ patterns map[string]*regexp.Regexp }

func NewRuleCompiler() *RuleCompiler { return &RuleCompiler{patterns: map[string]*regexp.Regexp{}} }
func (c *RuleCompiler) Compile(rule domain.MatchRule) (domain.MatchRule, error) {
	if err := domain.ValidateRule(rule); err != nil {
		return domain.MatchRule{}, err
	}
	if rule.Path != "" && !strings.HasSuffix(rule.Path, "*") {
		if _, ok := c.patterns[rule.Path]; !ok {
			re, e := regexp.Compile("^" + regexp.QuoteMeta(rule.Path) + "$")
			if e != nil {
				return domain.MatchRule{}, fmt.Errorf("compile path: %w", e)
			}
			c.patterns[rule.Path] = re
		}
	}
	rule.Method = domain.NormalizeMethod(rule.Method)
	rule.Path = domain.NormalizePath(rule.Path)
	return rule, nil
}
func (c *RuleCompiler) Match(rule domain.MatchRule, path string) bool {
	if strings.HasSuffix(rule.Path, "*") {
		return strings.HasPrefix(path, strings.TrimSuffix(rule.Path, "*"))
	}
	if re, ok := c.patterns[rule.Path]; ok {
		return re.MatchString(path)
	}
	return rule.Path == path
}
func RuleSpecificity(r domain.MatchRule) int {
	n := 0
	if r.ServiceID != "" {
		n += 10
	}
	if r.Environment != "" {
		n += 4
	}
	if r.Method != "" {
		n += 3
	}
	if r.Path != "" {
		n += 2
	}
	if r.SourceTag != "" {
		n += 1
	}
	return n
}
func SortBySpecificity(p []domain.Policy) []domain.Policy {
	for i := range p {
		for j := i + 1; j < len(p); j++ {
			if RuleSpecificity(p[j].Rule) > RuleSpecificity(p[i].Rule) {
				p[i], p[j] = p[j], p[i]
			}
		}
	}
	return p
}
