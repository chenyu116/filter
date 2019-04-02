package filter

import (
	"strconv"
	"strings"
)

type filterRule struct {
	Name          string
	In            string
	OutString     string
	OutInt        int
	Type          filterType
	Required      bool
	DefaultString string
	DefaultInt    int
}

type filterType string

type filter struct {
	rules []filterRule
}

func (f filter) AddRule(rule ...filterRule) {
	f.rules = append(f.rules, rule...)
}

func (f filter) GetString(key string) string {
	for _, v := range f.rules {
		if v.Name == key {
			return v.OutString
		}
	}
	return ""
}

func (f filter) GetInt(key string) int {
	for _, v := range f.rules {
		if v.Name == key {
			return v.OutInt
		}
	}
	return -1
}

func (f filter) Valid() bool {
	var r int
	var err error
	for _, v := range f.rules {
		v.In = strings.Trim(v.In, " ")
		if v.In == "" {
			if v.Required && v.DefaultString == "nil" {
				return false
			}
			v.OutString = v.DefaultString
			continue
		}
		if v.Type == FilterTypeUint() {
			r, err = strconv.Atoi(v.In)
			if err != nil {
				if v.Required && v.DefaultInt == -1 {
					return false
				} else {
					v.OutInt = v.DefaultInt
					continue
				}
			}
			if r < 0 {
				return false
			}
			v.OutInt = r
		} else if v.Type == FilterTypeString() {
			v.OutString = v.In
		} else {
			return false
		}
	}
	return true
}

func NewFilter() filter {
	return filter{}
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
		In:         value,
		Type:       FilterTypeUint(),
		DefaultInt: df,
		Required:   required,
	}
}

func RuleString(name, value string, required bool, df string) filterRule {
	return filterRule{
		Name:          name,
		In:            value,
		Type:          FilterTypeString(),
		DefaultString: df,
		Required:      required,
	}
}
