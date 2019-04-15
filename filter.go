package filter

import (
	"crypto/md5"
	"encoding/hex"
	"html"
	"sort"
	"strconv"
	"strings"
)

type filterRule struct {
	Name          string
	Value         string
	Type          filterType
	Required      bool
	DefaultString string
	DefaultInt    int
}

type filterType string

func GetString(rules []filterRule, key string, escape bool) string {
	for k := range rules {
		if rules[k].Name == key {
			if rules[k].Value == "" {
				return rules[k].DefaultString
			}
			if escape {
				return html.EscapeString(rules[k].Value)
			}
		}
	}
	return ""
}

func GetInt(rules []filterRule, key string) int {
	for k := range rules {
		if rules[k].Name == key {
			if rules[k].Value == "" {
				return rules[k].DefaultInt
			}
			rn, err := strconv.Atoi(rules[k].Value)
			if err != nil || rn < 0 {
				return rules[k].DefaultInt
			}
			return rn
		}
	}
	return -1
}

func Valid(rules ...filterRule) ([]filterRule, bool) {
	var r int
	var err error
	for k := range rules {
		if rules[k].Type == FilterTypeUint() {
			if rules[k].Required {
				if rules[k].Value == "" && rules[k].DefaultInt == -1 {
					return nil, false
				} else {
					continue
				}
			}
			r, err = strconv.Atoi(rules[k].Value)
			if err != nil {
				if rules[k].Required && rules[k].DefaultInt == -1 {
					return nil, false
				} else {
					continue
				}
			}
			if r < 0 {
				return nil, false
			}
			continue
		} else if rules[k].Type == FilterTypeString() {
			if rules[k].Required {
				if rules[k].Value == "" && rules[k].DefaultString == "" {
					return nil, false
				}
			}
			continue
		} else {
			return nil, false
		}
	}
	return rules, true
}

func ValidToken(r []filterRule, KEY, token string) bool {
	if token == "" {
		return true
	}

	var keys []string
	for k := range r {
		if r[k].Required {
			if r[k].Name == "app" && r[k].Value == "mobile" {
				return true
			}
			keys = append(keys, r[k].Name)
		}
	}
	sort.Strings(keys)
	var tokenString string
	for k := range keys {
		for rk := range r {
			if r[rk].Name == keys[k] {
				tokenString += r[rk].Value
				break
			}
		}
	}
	tokenBytes := md5.Sum([]byte(tokenString + KEY))
	if strings.Compare(token, hex.EncodeToString(tokenBytes[:])) != 0 {
		return false
	}
	return true
}

func FilterTypeUint() filterType {
	return "uint"
}

func FilterTypeString() filterType {
	return "string"
}

func RuleUint(name, value string, required bool, df int) filterRule {
	return filterRule{
		Name:       name,
		Value:      value,
		Type:       FilterTypeUint(),
		DefaultInt: df,
		Required:   required,
	}
}

func RuleString(name, value string, required bool, df string) filterRule {
	return filterRule{
		Name:          name,
		Value:         value,
		Type:          FilterTypeString(),
		DefaultString: df,
		Required:      required,
	}
}
