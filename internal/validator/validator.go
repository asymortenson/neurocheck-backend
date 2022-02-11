package validator

import (
	"net/url"
	"regexp"
	"strconv"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

func ValidateQuery(values url.Values, query []string) (string, bool) {
	if !values.Has("domain") && !values.Has("owner_id") {
		return "domain, " + "owner_id", false
	}

	for _, key := range query {
		if !values.Has(key) {
			return key, false
		}
	}

	count, _ := strconv.Atoi(values.Get("count"))

	if count > 10 {
		return "count cannot be greater than 10", false
	}

	if values.Has("type") {
		if values.Get("type") != "result" && values.Get("type") != "toxic" {
			return values.Get("type") + " is bad type", false
		}
		return "", true
	}

	return "", true
}
