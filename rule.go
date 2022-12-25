package main

import (
	"encoding/json"
	"regexp"
)

type RegexRule struct {
	rawPattern string
	rule       *regexp.Regexp
}

func (r *RegexRule) UnmarshalJSON(data []byte) error {
	var regexPattern string
	if err := json.Unmarshal(data, &regexPattern); err != nil {
		return err
	}

	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return err
	}

	r.rawPattern = regexPattern
	r.rule = regex
	return nil
}

func (r *RegexRule) IsMatching(str string) bool {
	return r.rule.MatchString(str)
}

func (r *RegexRule) Replace(str string, repl string) string {
	return r.rule.ReplaceAllString(str, repl)
}

func isMatchingAnyRule(rules []RegexRule, str string) bool {
	for _, rule := range rules {
		if rule.IsMatching(str) {
			return true
		}
	}

	return false
}
