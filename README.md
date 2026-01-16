# go-lua-validation

A GopherLua module for data validation in Lua scripts.

## Installation

```bash
go get github.com/kerimovok/go-lua-validation
```

## Usage

### In Go

```go
import (
    lua "github.com/yuin/gopher-lua"
    validation "github.com/kerimovok/go-lua-validation"
)

L := lua.NewState()
defer L.Close()

// Preload the validation module
L.PreloadModule("validation", validation.Loader)

// Now Lua scripts can use require("validation")
L.DoString(`
    local validation = require("validation")
    if validation.validate_email("user@example.com") then
        print("Valid email")
    end
`)
```

### In Lua

```lua
local validation = require("validation")

-- Type checking
if validation.is_string(value) then
    print("It's a string")
end

if validation.is_number(value) then
    print("It's a number")
end

if validation.is_table(value) then
    print("It's a table")
end

-- Empty check
if validation.is_empty(value) then
    error("Value cannot be empty")
end

-- Email validation
if not validation.validate_email(email) then
    error("Invalid email address")
end

-- URL validation
if not validation.validate_url(url) then
    error("Invalid URL")
end

-- Regex validation
local isValid, err = validation.validate_regex("abc123", "^[a-z]+[0-9]+$")
if err then
    error("Invalid regex pattern: " .. err)
elseif not isValid then
    error("String does not match pattern")
end

-- Length validation
if not validation.min_length(password, 8) then
    error("Password must be at least 8 characters")
end

if not validation.max_length(username, 50) then
    error("Username must be at most 50 characters")
end

-- Range validation
if not validation.in_range(age, 18, 120) then
    error("Age must be between 18 and 120")
end
```

## Functions

### Type Checking

#### `validation.is_string(value)`

Checks if a value is a string.

- **Parameters:**
  - `value`: Value to check
- **Returns:**
  - `boolean`: `true` if string, `false` otherwise

#### `validation.is_number(value)`

Checks if a value is a number.

- **Parameters:**
  - `value`: Value to check
- **Returns:**
  - `boolean`: `true` if number, `false` otherwise

#### `validation.is_table(value)`

Checks if a value is a table.

- **Parameters:**
  - `value`: Value to check
- **Returns:**
  - `boolean`: `true` if table, `false` otherwise

#### `validation.is_boolean(value)`

Checks if a value is a boolean.

- **Parameters:**
  - `value`: Value to check
- **Returns:**
  - `boolean`: `true` if boolean, `false` otherwise

#### `validation.is_nil(value)`

Checks if a value is nil.

- **Parameters:**
  - `value`: Value to check
- **Returns:**
  - `boolean`: `true` if nil, `false` otherwise

### Value Validation

#### `validation.is_empty(value)`

Checks if a value is empty (nil, empty string, or empty table).

- **Parameters:**
  - `value`: Value to check
- **Returns:**
  - `boolean`: `true` if empty, `false` otherwise

### Format Validation

#### `validation.validate_email(email)`

Validates an email address.

- **Parameters:**
  - `email` (string): Email address to validate
- **Returns:**
  - `boolean`: `true` if valid email, `false` otherwise

#### `validation.validate_url(url)`

Validates a URL.

- **Parameters:**
  - `url` (string): URL to validate
- **Returns:**
  - `boolean`: `true` if valid URL, `false` otherwise

#### `validation.validate_regex(str, pattern)`

Validates a string against a regex pattern.

- **Parameters:**
  - `str` (string): String to validate
  - `pattern` (string): Regex pattern
- **Returns:**
  - `boolean`: `true` if matches, `false` otherwise (or `nil` if regex pattern is invalid)
  - `string` (error): Error message if regex pattern is invalid (only returned on error)

### Length Validation

#### `validation.min_length(str, min)`

Checks if a string has minimum length.

- **Parameters:**
  - `str` (string): String to check
  - `min` (number): Minimum length
- **Returns:**
  - `boolean`: `true` if length >= min, `false` otherwise

#### `validation.max_length(str, max)`

Checks if a string has maximum length.

- **Parameters:**
  - `str` (string): String to check
  - `max` (number): Maximum length
- **Returns:**
  - `boolean`: `true` if length <= max, `false` otherwise

### Range Validation

#### `validation.in_range(num, min, max)`

Checks if a number is within a range.

- **Parameters:**
  - `num` (number): Number to check
  - `min` (number): Minimum value
  - `max` (number): Maximum value
- **Returns:**
  - `boolean`: `true` if min <= num <= max, `false` otherwise

## Notes

- Email validation uses Go's `net/mail` package
- URL validation uses Go's `net/url` package
- Regex patterns use Go's regex syntax (RE2)
- All validation functions are safe and do not throw errors (except `validate_regex` which may return an error for invalid patterns)
