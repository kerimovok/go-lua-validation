package validation

import (
	"net/mail"
	"net/url"
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

// Loader loads the validation module
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"is_empty":       isEmpty,
	"is_string":      isString,
	"is_number":      isNumber,
	"is_table":       isTable,
	"is_boolean":     isBoolean,
	"is_nil":         isNil,
	"validate_email": validateEmail,
	"validate_url":   validateURL,
	"validate_regex": validateRegex,
	"min_length":     minLength,
	"max_length":     maxLength,
	"in_range":       inRange,
}

// isEmpty checks if a value is nil, empty string, or empty table
// Usage: validation.is_empty(value) -> boolean
func isEmpty(L *lua.LState) int {
	value := L.CheckAny(1)

	if value == lua.LNil {
		L.Push(lua.LBool(true))
		return 1
	}

	if str, ok := value.(lua.LString); ok {
		L.Push(lua.LBool(string(str) == ""))
		return 1
	}

	if tbl, ok := value.(*lua.LTable); ok {
		count := 0
		tbl.ForEach(func(_, _ lua.LValue) {
			count++
		})
		L.Push(lua.LBool(count == 0))
		return 1
	}

	L.Push(lua.LBool(false))
	return 1
}

// isString checks if a value is a string
// Usage: validation.is_string(value) -> boolean
func isString(L *lua.LState) int {
	value := L.CheckAny(1)
	_, ok := value.(lua.LString)
	L.Push(lua.LBool(ok))
	return 1
}

// isNumber checks if a value is a number
// Usage: validation.is_number(value) -> boolean
func isNumber(L *lua.LState) int {
	value := L.CheckAny(1)
	_, ok := value.(lua.LNumber)
	L.Push(lua.LBool(ok))
	return 1
}

// isTable checks if a value is a table
// Usage: validation.is_table(value) -> boolean
func isTable(L *lua.LState) int {
	value := L.CheckAny(1)
	_, ok := value.(*lua.LTable)
	L.Push(lua.LBool(ok))
	return 1
}

// isBoolean checks if a value is a boolean
// Usage: validation.is_boolean(value) -> boolean
func isBoolean(L *lua.LState) int {
	value := L.CheckAny(1)
	_, ok := value.(lua.LBool)
	L.Push(lua.LBool(ok))
	return 1
}

// isNil checks if a value is nil
// Usage: validation.is_nil(value) -> boolean
func isNil(L *lua.LState) int {
	value := L.CheckAny(1)
	L.Push(lua.LBool(value == lua.LNil))
	return 1
}

// validateEmail validates an email address
// Usage: validation.validate_email(email) -> boolean
func validateEmail(L *lua.LState) int {
	email := L.CheckString(1)
	_, err := mail.ParseAddress(email)
	L.Push(lua.LBool(err == nil))
	return 1
}

// validateURL validates a URL
// Usage: validation.validate_url(url) -> boolean
func validateURL(L *lua.LState) int {
	urlStr := L.CheckString(1)
	_, err := url.ParseRequestURI(urlStr)
	L.Push(lua.LBool(err == nil))
	return 1
}

// validateRegex validates a string against a regex pattern
// Usage: validation.validate_regex(str, pattern) -> boolean, error?
func validateRegex(L *lua.LState) int {
	str := L.CheckString(1)
	pattern := L.CheckString(2)

	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(re.MatchString(str)))
	return 1
}

// minLength checks if a string has minimum length
// Usage: validation.min_length(str, min) -> boolean
func minLength(L *lua.LState) int {
	str := L.CheckString(1)
	min := L.CheckInt(2)
	L.Push(lua.LBool(len(str) >= min))
	return 1
}

// maxLength checks if a string has maximum length
// Usage: validation.max_length(str, max) -> boolean
func maxLength(L *lua.LState) int {
	str := L.CheckString(1)
	max := L.CheckInt(2)
	L.Push(lua.LBool(len(str) <= max))
	return 1
}

// inRange checks if a number is within a range
// Usage: validation.in_range(num, min, max) -> boolean
func inRange(L *lua.LState) int {
	num := L.CheckNumber(1)
	min := L.CheckNumber(2)
	max := L.CheckNumber(3)
	L.Push(lua.LBool(num >= min && num <= max))
	return 1
}
