# pform

[![Go Reference](https://pkg.go.dev/badge/github.com/ElOrlis/pform.svg)](https://pkg.go.dev/github.com/ElOrlis/pform)
[![Code Validation](https://github.com/ElOrlis/pform/actions/workflows/code-validation.yml/badge.svg)](https://github.com/ElOrlis/pform/actions/workflows/code-validation.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ElOrlis/pform)](https://goreportcard.com/report/github.com/ElOrlis/pform)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight Go library for decoding `application/x-www-form-urlencoded` HTTP request bodies into Go structs or maps. Uses struct tags for field mapping and supports custom type unmarshalling.

## Features

- Decode form data directly from `http.Request` or `url.Values`
- Support for all Go primitive types (strings, integers, floats, booleans, complex numbers)
- Struct tag-based field mapping with flexible options
- Custom type unmarshalling via the `Unmarshaller` interface
- Map decoding support
- Detailed error types for parse failures and missing required fields
- Zero external dependencies (except for testing)

## Installation

```bash
go get github.com/ElOrlis/pform
```

## Quick Start

### Decoding from HTTP Request

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/ElOrlis/pform"
)

type LoginForm struct {
    Username string `form:"username,required"`
    Password string `form:"password,required"`
    Remember bool   `form:"remember,omitempty"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    decoder, err := pform.NewFormUrlDecoder(r)
    if err != nil {
        http.Error(w, "Invalid content type", http.StatusBadRequest)
        return
    }

    var form LoginForm
    if err := decoder.Decode(&form); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Printf("Login attempt: %s\n", form.Username)
}
```

### Decoding from url.Values

```go
package main

import (
    "fmt"
    "net/url"

    "github.com/ElOrlis/pform"
)

type UserProfile struct {
    Name    string  `form:"name"`
    Age     int     `form:"age"`
    Email   string  `form:"email,omitempty"`
    Balance float64 `form:"balance"`
}

func main() {
    values := url.Values{
        "name":    {"John Doe"},
        "age":     {"30"},
        "balance": {"1234.56"},
    }

    var profile UserProfile
    decoder := pform.NewDecoder(values)
    if err := decoder.Decode(&profile); err != nil {
        panic(err)
    }

    fmt.Printf("Name: %s, Age: %d, Balance: %.2f\n",
        profile.Name, profile.Age, profile.Balance)
}
```

## Struct Tags

The `form` struct tag controls how fields are mapped and validated:

```go
type Example struct {
    // Maps to form field "field_name"
    Field1 string `form:"field_name"`

    // Skip empty values (no error if missing)
    Field2 string `form:"field_name,omitempty"`

    // Return error if field is missing
    Field3 string `form:"field_name,required"`

    // Ignore this field completely
    Field4 string `form:"-"`

    // Fields without tags are ignored
    Field5 string
}
```

### Tag Options

| Option | Description |
|--------|-------------|
| `form:"name"` | Maps the struct field to form field "name" |
| `form:"name,omitempty"` | Skip if the form field is empty or missing |
| `form:"name,required"` | Return `RequiredFieldError` if the field is missing |
| `form:"-"` | Ignore this field during decoding |

## Supported Types

### Primitive Types

| Type | Example Values |
|------|----------------|
| `string` | `"hello"`, `"John Doe"` |
| `int`, `int8`, `int16`, `int32`, `int64` | `"42"`, `"-100"` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `"42"`, `"255"` |
| `float32`, `float64` | `"3.14"`, `"1.23e-4"` |
| `complex64`, `complex128` | `"(1+2i)"`, `"3.14"`, `"5i"` |
| `bool` | `"true"`, `"false"`, `"1"`, `"0"`, `"t"`, `"f"` |

### Map Types

Decode form data directly into a map:

```go
values := url.Values{
    "key1": {"value1"},
    "key2": {"value2"},
}

var data map[string]string
decoder := pform.NewDecoder(values)
err := decoder.Decode(&data)
// data = map[string]string{"key1": "value1", "key2": "value2"}
```

Maps support any value type that can be parsed:

```go
var counts map[string]int
var scores map[string]float64
```

**Note:** Map keys must be of type `string`.

## Custom Types

Implement the `Unmarshaller` interface to handle custom type parsing:

```go
type Unmarshaller interface {
    UnmarshalValue(string) error
}
```

### Example: Custom Date Type

```go
type Date struct {
    Year  int
    Month int
    Day   int
}

func (d *Date) UnmarshalValue(v string) error {
    t, err := time.Parse("2006-01-02", v)
    if err != nil {
        return err
    }
    d.Year = t.Year()
    d.Month = int(t.Month())
    d.Day = t.Day()
    return nil
}

type Event struct {
    Name string `form:"name"`
    Date Date   `form:"date"`
}
```

### Example: Custom Enum Type

```go
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
)

func (s *Status) UnmarshalValue(v string) error {
    switch strings.ToLower(v) {
    case "pending":
        *s = StatusPending
    case "active":
        *s = StatusActive
    case "inactive":
        *s = StatusInactive
    default:
        return fmt.Errorf("invalid status: %s", v)
    }
    return nil
}
```

## Error Handling

The library provides two custom error types for detailed error information:

### ParseError

Returned when a value cannot be converted to the target type:

```go
var parseErr pform.ParseError
if errors.As(err, &parseErr) {
    fmt.Printf("Field: %s\n", parseErr.Field)   // Field name
    fmt.Printf("Type: %s\n", parseErr.Type)     // Expected type
    fmt.Printf("Value: %s\n", parseErr.Value)   // Provided value
    fmt.Printf("Cause: %v\n", parseErr.Err)     // Underlying error
}
```

### RequiredFieldError

Returned when a required field is missing:

```go
var reqErr pform.RequiredFieldError
if errors.As(err, &reqErr) {
    fmt.Printf("Missing required field: %s\n", reqErr.Error())
}
```

### Example Error Handling

```go
func handleForm(r *http.Request) error {
    decoder, err := pform.NewFormUrlDecoder(r)
    if err != nil {
        return fmt.Errorf("invalid request: %w", err)
    }

    var form MyForm
    if err := decoder.Decode(&form); err != nil {
        var parseErr pform.ParseError
        var reqErr pform.RequiredFieldError

        switch {
        case errors.As(err, &reqErr):
            return fmt.Errorf("missing required field: %w", err)
        case errors.As(err, &parseErr):
            return fmt.Errorf("invalid value for %s: %w", parseErr.Field, err)
        default:
            return fmt.Errorf("decode error: %w", err)
        }
    }

    return nil
}
```

## API Reference

### Functions

#### NewFormUrlDecoder

```go
func NewFormUrlDecoder(r *http.Request) (XFormUrlEncoded, error)
```

Creates a decoder from an HTTP request. Validates that the Content-Type is `application/x-www-form-urlencoded`.

#### NewDecoder

```go
func NewDecoder(form url.Values) Decoder
```

Creates a decoder from `url.Values` directly.

### Types

#### Decoder

```go
type Decoder struct {
    // contains filtered or unexported fields
}

func (d Decoder) Decode(dest any) error
```

Decodes form values into the destination struct or map. The destination must be a pointer to a struct or map.

#### XFormUrlEncoded

```go
type XFormUrlEncoded interface {
    Decode(dest any) error
}
```

Interface implemented by the decoder.

#### Unmarshaller

```go
type Unmarshaller interface {
    UnmarshalValue(string) error
}
```

Implement this interface to provide custom parsing logic for your types.

## Content-Type Handling

The `NewFormUrlDecoder` function validates the HTTP request's Content-Type header:

- Must be `application/x-www-form-urlencoded`
- Accepts charset parameters (e.g., `application/x-www-form-urlencoded; charset=utf-8`)
- Case-insensitive matching
- Trims whitespace

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new functionality
- Ensure all tests pass with `go test ./...`
- Run `go vet ./...` and `staticcheck ./...` before submitting
- Format code with `gofmt`

## LLM Usage

This project uses LLM assistance exclusively for:
- **Testing** - Generating test cases and test coverage
- **CI/CD Workflows** - Creating and maintaining GitHub Actions workflows
- **Documentation** - Writing and improving documentation

The core library code is human-written without LLM assistance.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [gorilla/schema](https://github.com/gorilla/schema) - Package for converting structs to and from form values
- [go-playground/form](https://github.com/go-playground/form) - Decodes url.Values into Go value(s) and Encodes Go value(s) into url.Values
