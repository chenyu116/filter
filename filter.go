package filter

import (
	"strconv"
	"strings"
)

type filterRule struct {
	Name     string
	Value    string
	Type     filterType
	Required bool
	Default  interface{}
}

type filterType string

type filter struct {
	rules  []filterRule
	data   map[string]interface{}
	errMsg string
}

func (f *filter) AddRule(rule ...filterRule) {
	f.rules = append(f.rules, rule...)
}

func (f *filter) Map() map[string]interface{} {
	return f.data
}

func (f *filter) Get(key string) interface{} {
	if data, ok := f.data[key]; ok {
		return data
	}
	return nil
}

func (f *filter) ErrorMsg() string {
	return f.errMsg
}

func (f *filter) Valid() bool {
	for _, v := range f.rules {
		v.Value = strings.Trim(v.Value, " ")
		if v.Value == "" {
			f.data[v.Name] = v.Default
			if v.Required && v.Default == nil {
				f.errMsg = v.Name + " is empty"
				return false
			}
			continue
		}
		if v.Type == FilterTypeUint() {
			r, err := strconv.Atoi(v.Value)
			if err != nil {
				if v.Default != nil {
					r = v.Default.(int)
				} else if v.Required {
					f.errMsg = v.Name + " is invalid"
					return false
				} else {
					f.data[v.Name] = nil
					continue
				}
			}
			if r < 0 {
				f.errMsg = v.Name + " is invalid"
				return false
			}
			f.data[v.Name] = r
		} else if v.Type == FilterTypeString() {
			f.data[v.Name] = v.Value
		} else {
			f.errMsg = "Unknow type " + string(v.Type)
			return false
		}
	}
	return true
}

func NewFilter() *filter {
	return &filter{
		data: make(map[string]interface{}),
	}
}

func FilterTypeUint() filterType {
	return "uint"
}

func FilterTypeString() filterType {
	return "string"
}

func Rule(name, value string, ft filterType, required bool, df interface{}) filterRule {
	return filterRule{
		Name:     name,
		Value:    value,
		Type:     ft,
		Default:  df,
		Required: required,
	}
}

func RuleUint(name, value string, required bool, df interface{}) filterRule {
	return filterRule{
		Name:     name,
		Value:    value,
		Type:     FilterTypeUint(),
		Default:  df,
		Required: required,
	}
}

func RuleString(name, value string, required bool, df interface{}) filterRule {
	return filterRule{
		Name:     name,
		Value:    value,
		Type:     FilterTypeString(),
		Default:  df,
		Required: required,
	}
}
