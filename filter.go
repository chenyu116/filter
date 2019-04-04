package filter

import (
	"crypto/md5"
	"encoding/hex"
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

func GetString(rules []filterRule, key string) string {
	for _, v := range rules {
		if v.Name == key {
			if v.Value == "" {
				return v.DefaultString
			}
			return v.Value
		}
	}
	return ""
}

func GetInt(rules []filterRule, key string) int {
	for _, v := range rules {
		if v.Name == key {
			v.Value = strings.Trim(v.Value, " ")
			if v.Value == "" {
				return v.DefaultInt
			}
			r, err := strconv.Atoi(v.Value)
			if err != nil || r < 0 {
				return v.DefaultInt
			}
			return r
		}
	}
	return -1
}

func Valid(rules ...filterRule) ([]filterRule, bool) {
	var r int
	var err error
	for _, v := range rules {
		v.Value = strings.Trim(v.Value, " ")
		if v.Value == "" {
			if v.Required && v.DefaultString == "" {
				return nil, false
			}
			continue
		}
		if v.Type == FilterTypeUint() {
			r, err = strconv.Atoi(v.Value)
			if err != nil {
				if v.Required && v.DefaultInt == -1 {
					return nil, false
				} else {
					continue
				}
			}
			if r < 0 {
				return nil, false
			}
			continue
		} else if v.Type == FilterTypeString() {
			continue
		} else {
			return nil, false
		}
	}
	return rules, true
}

func ValidToken(r []filterRule, KEY, token string) bool {
	var keys []string
	for k := range r {
		if r[k].Required {
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
	_token := hex.EncodeToString(tokenBytes[:])
	if strings.Compare(token, _token) != 0 {
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
