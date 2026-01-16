package validation

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLoader(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	err := L.DoString(`
		local validation = require("validation")
		if validation == nil then
			error("validation module is nil")
		end
	`)
	if err != nil {
		t.Fatalf("Failed to load validation module: %v", err)
	}
}

func TestIsString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"string", `"hello"`, true},
		{"number", "42", false},
		{"table", "{}", false},
		{"boolean", "true", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := `
				local validation = require("validation")
				return validation.is_string(` + tt.value + `)
			`

			err := L.DoString(script)
			if err != nil {
				t.Fatalf("IsString test failed: %v", err)
			}

			result := L.Get(-1).(lua.LBool)
			if bool(result) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsNumber(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	script := `
		local validation = require("validation")
		return validation.is_number(42), validation.is_number("42")
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("IsNumber test failed: %v", err)
	}

	result1 := L.Get(-2).(lua.LBool)
	result2 := L.Get(-1).(lua.LBool)

	if !bool(result1) {
		t.Error("Expected true for number 42")
	}
	if bool(result2) {
		t.Error("Expected false for string '42'")
	}
}

func TestIsTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	script := `
		local validation = require("validation")
		return validation.is_table({}), validation.is_table("not a table")
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("IsTable test failed: %v", err)
	}

	result1 := L.Get(-2).(lua.LBool)
	result2 := L.Get(-1).(lua.LBool)

	if !bool(result1) {
		t.Error("Expected true for table")
	}
	if bool(result2) {
		t.Error("Expected false for string")
	}
}

func TestIsEmpty(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	script := `
		local validation = require("validation")
		return validation.is_empty(nil), validation.is_empty(""), validation.is_empty({}), validation.is_empty("hello")
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("IsEmpty test failed: %v", err)
	}

	nilResult := L.Get(-4).(lua.LBool)
	emptyStr := L.Get(-3).(lua.LBool)
	emptyTable := L.Get(-2).(lua.LBool)
	nonEmpty := L.Get(-1).(lua.LBool)

	if !bool(nilResult) {
		t.Error("Expected true for nil")
	}
	if !bool(emptyStr) {
		t.Error("Expected true for empty string")
	}
	if !bool(emptyTable) {
		t.Error("Expected true for empty table")
	}
	if bool(nonEmpty) {
		t.Error("Expected false for non-empty string")
	}
}

func TestValidateEmail(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"valid email", "user@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"invalid email", "not-an-email", false},
		{"invalid email no domain", "user@", false},
		{"invalid email no @", "userexample.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := `
				local validation = require("validation")
				return validation.validate_email("` + tt.email + `")
			`

			err := L.DoString(script)
			if err != nil {
				t.Fatalf("ValidateEmail test failed: %v", err)
			}

			result := L.Get(-1).(lua.LBool)
			if bool(result) != tt.expected {
				t.Errorf("Expected %v for %s, got %v", tt.expected, tt.email, result)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"valid http URL", "http://example.com", true},
		{"valid https URL", "https://example.com", true},
		{"valid URL with path", "https://example.com/path", true},
		{"invalid URL", "not-a-url", false},
		{"invalid URL no scheme", "example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := `
				local validation = require("validation")
				return validation.validate_url("` + tt.url + `")
			`

			err := L.DoString(script)
			if err != nil {
				t.Fatalf("ValidateURL test failed: %v", err)
			}

			result := L.Get(-1).(lua.LBool)
			if bool(result) != tt.expected {
				t.Errorf("Expected %v for %s, got %v", tt.expected, tt.url, result)
			}
		})
	}
}

func TestValidateRegex(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	script := `
		local validation = require("validation")
		local isValid, err = validation.validate_regex("abc123", "^[a-z]+[0-9]+$")
		if err then
			error("Regex validation failed: " .. err)
		end
		return isValid
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("ValidateRegex test failed: %v", err)
	}

	result := L.Get(-1).(lua.LBool)
	if !bool(result) {
		t.Error("Expected true for matching regex")
	}
}

func TestValidateRegexInvalidPattern(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	script := `
		local validation = require("validation")
		local isValid, err = validation.validate_regex("test", "[invalid")
		if err == nil then
			error("Expected error for invalid regex pattern")
		end
		return isValid, err
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("ValidateRegex invalid pattern test failed: %v", err)
	}

	result := L.Get(-2)
	errVal := L.Get(-1)
	// When regex is invalid, isValid should be false (not nil) and err should be set
	if result == lua.LNil || errVal == lua.LNil {
		t.Errorf("Expected isValid=false and error string, got isValid=%v, err=%v", result, errVal)
	}
}

func TestMinLength(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	script := `
		local validation = require("validation")
		return validation.min_length("hello", 3), validation.min_length("hi", 3)
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("MinLength test failed: %v", err)
	}

	result1 := L.Get(-2).(lua.LBool)
	result2 := L.Get(-1).(lua.LBool)

	if !bool(result1) {
		t.Error("Expected true for 'hello' with min length 3")
	}
	if bool(result2) {
		t.Error("Expected false for 'hi' with min length 3")
	}
}

func TestMaxLength(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	script := `
		local validation = require("validation")
		return validation.max_length("hello", 10), validation.max_length("hello", 3)
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("MaxLength test failed: %v", err)
	}

	result1 := L.Get(-2).(lua.LBool)
	result2 := L.Get(-1).(lua.LBool)

	if !bool(result1) {
		t.Error("Expected true for 'hello' with max length 10")
	}
	if bool(result2) {
		t.Error("Expected false for 'hello' with max length 3")
	}
}

func TestInRange(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("validation", Loader)

	script := `
		local validation = require("validation")
		return validation.in_range(5, 1, 10), validation.in_range(15, 1, 10), validation.in_range(0, 1, 10)
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("InRange test failed: %v", err)
	}

	result1 := L.Get(-3).(lua.LBool)
	result2 := L.Get(-2).(lua.LBool)
	result3 := L.Get(-1).(lua.LBool)

	if !bool(result1) {
		t.Error("Expected true for 5 in range [1, 10]")
	}
	if bool(result2) {
		t.Error("Expected false for 15 in range [1, 10]")
	}
	if bool(result3) {
		t.Error("Expected false for 0 in range [1, 10]")
	}
}
