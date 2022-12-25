package main

import (
	"net/url"
)

type Rules struct {
	Providers map[string]Provider `json:"providers"`
}

func (r *Rules) CleanUrl(url string, allowReferralMarketing bool) (string, error) {
	var err error

	for _, provider := range r.Providers {
		if provider.MatchesUrl(url) {
			url, err = provider.CleanUrl(url, allowReferralMarketing)
			if err != nil {
				return "", err
			}
		}
	}
	return url, nil
}

type Provider struct {
	UrlPattern        RegexRule   `json:"urlPattern"`
	CompleteProvider  bool        `json:"completeProvider"`
	Rules             []RegexRule `json:"rules"`
	ReferralMarketing []RegexRule `json:"referralMarketing"`
	Exceptions        []RegexRule `json:"exceptions"`
	RawRules          []RegexRule `json:"rawRules"`
	Redirections      []RegexRule `json:"redirections"`
	ForceRedirection  bool        `json:"forceRedirection"`
}

func (p *Provider) MatchesUrl(url string) bool {
	if !p.UrlPattern.IsMatching(url) {
		return false
	}

	for _, exception := range p.Exceptions {
		if exception.IsMatching(url) {
			return false
		}
	}

	return true
}

func (p *Provider) CleanUrl(urlStr string, allowReferralMarketing bool) (string, error) {
	// apply raw rules
	for _, rawRule := range p.RawRules {
		urlStr = rawRule.Replace(urlStr, "")
	}

	// collect rules and additionally referral marketing rules
	rules := make([]RegexRule, len(p.Rules))
	copy(rules, p.Rules)

	if !allowReferralMarketing {
		rules = append(rules, p.ReferralMarketing...)
	}

	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	if parsedUrl.Scheme == "" {
		return urlStr, nil
	}

	query := parsedUrl.Query()
	queryParams := mapKeys(query)
	for _, queryParamName := range queryParams {
		if isMatchingAnyRule(rules, queryParamName) {
			query.Del(queryParamName)
		}
	}

	parsedUrl.RawQuery = query.Encode()
	return parsedUrl.String(), nil
}
