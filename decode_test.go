package pform_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"pform"
)

// Custom type with pointer receiver
type CustomType struct {
	Value string
}

func (c *CustomType) UnmarshalValue(v string) error {
	c.Value = "custom:" + v
	return nil
}

// Custom type with value receiver (tests pointer-to-type handling)
type CustomValueType struct {
	Value string
}

func (c CustomValueType) UnmarshalValue(v string) error {
	c.Value = "value:" + v
	return nil
}

// Custom type that returns errors
type CustomErrorType struct {
	Value string
}

func (c *CustomErrorType) UnmarshalValue(v string) error {
	if v == "error" {
		return errors.New("custom error")
	}
	c.Value = v
	return nil
}

// Test basic struct decoding with all integer types
func TestDecoder_IntegerTypes(t *testing.T) {
	type TestStruct struct {
		Int   int   `form:"int"`
		Int8  int8  `form:"int8"`
		Int16 int16 `form:"int16"`
		Int32 int32 `form:"int32"`
		Int64 int64 `form:"int64"`
	}

	tests := []struct {
		name    string
		values  url.Values
		want    TestStruct
		wantErr bool
	}{
		{
			name: "valid positive integers",
			values: url.Values{
				"int":   {"42"},
				"int8":  {"127"},
				"int16": {"32767"},
				"int32": {"2147483647"},
				"int64": {"9223372036854775807"},
			},
			want: TestStruct{
				Int:   42,
				Int8:  127,
				Int16: 32767,
				Int32: 2147483647,
				Int64: 9223372036854775807,
			},
			wantErr: false,
		},
		{
			name: "valid negative integers",
			values: url.Values{
				"int":   {"-42"},
				"int8":  {"-128"},
				"int16": {"-32768"},
				"int32": {"-2147483648"},
				"int64": {"-9223372036854775808"},
			},
			want: TestStruct{
				Int:   -42,
				Int8:  -128,
				Int16: -32768,
				Int32: -2147483648,
				Int64: -9223372036854775808,
			},
			wantErr: false,
		},
		{
			name: "invalid integer",
			values: url.Values{
				"int": {"not_a_number"},
			},
			wantErr: true,
		},
		{
			name: "int8 overflow",
			values: url.Values{
				"int8": {"128"},
			},
			wantErr: true,
		},
		{
			name: "int16 overflow",
			values: url.Values{
				"int16": {"32768"},
			},
			wantErr: true,
		},
		{
			name: "int32 overflow",
			values: url.Values{
				"int32": {"2147483648"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestStruct
			dec := pform.NewDecoder(tt.values)
			err := dec.Decode(&got)

			if tt.wantErr {
				assert.Error(t, err)
				var parseErr pform.ParseError
				assert.True(t, errors.As(err, &parseErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// Test unsigned integer types
func TestDecoder_UnsignedIntegerTypes(t *testing.T) {
	type TestStruct struct {
		Uint   uint   `form:"uint"`
		Uint8  uint8  `form:"uint8"`
		Uint16 uint16 `form:"uint16"`
		Uint32 uint32 `form:"uint32"`
		Uint64 uint64 `form:"uint64"`
	}

	tests := []struct {
		name    string
		values  url.Values
		want    TestStruct
		wantErr bool
	}{
		{
			name: "valid unsigned integers",
			values: url.Values{
				"uint":   {"42"},
				"uint8":  {"255"},
				"uint16": {"65535"},
				"uint32": {"4294967295"},
				"uint64": {"18446744073709551615"},
			},
			want: TestStruct{
				Uint:   42,
				Uint8:  255,
				Uint16: 65535,
				Uint32: 4294967295,
				Uint64: 18446744073709551615,
			},
			wantErr: false,
		},
		{
			name: "zero values",
			values: url.Values{
				"uint":   {"0"},
				"uint8":  {"0"},
				"uint16": {"0"},
				"uint32": {"0"},
				"uint64": {"0"},
			},
			want: TestStruct{
				Uint:   0,
				Uint8:  0,
				Uint16: 0,
				Uint32: 0,
				Uint64: 0,
			},
			wantErr: false,
		},
		{
			name: "negative value error",
			values: url.Values{
				"uint": {"-1"},
			},
			wantErr: true,
		},
		{
			name: "uint8 overflow",
			values: url.Values{
				"uint8": {"256"},
			},
			wantErr: true,
		},
		{
			name: "uint16 overflow",
			values: url.Values{
				"uint16": {"65536"},
			},
			wantErr: true,
		},
		{
			name: "invalid unsigned integer",
			values: url.Values{
				"uint": {"not_a_number"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestStruct
			dec := pform.NewDecoder(tt.values)
			err := dec.Decode(&got)

			if tt.wantErr {
				assert.Error(t, err)
				var parseErr pform.ParseError
				assert.True(t, errors.As(err, &parseErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// Test float types
func TestDecoder_FloatTypes(t *testing.T) {
	type TestStruct struct {
		Float32 float32 `form:"float32"`
		Float64 float64 `form:"float64"`
	}

	tests := []struct {
		name    string
		values  url.Values
		want    TestStruct
		wantErr bool
	}{
		{
			name: "valid floats",
			values: url.Values{
				"float32": {"3.14"},
				"float64": {"2.718281828459045"},
			},
			want: TestStruct{
				Float32: 3.14,
				Float64: 2.718281828459045,
			},
			wantErr: false,
		},
		{
			name: "scientific notation",
			values: url.Values{
				"float32": {"1.23e-4"},
				"float64": {"5.67e8"},
			},
			want: TestStruct{
				Float32: 1.23e-4,
				Float64: 5.67e8,
			},
			wantErr: false,
		},
		{
			name: "negative floats",
			values: url.Values{
				"float32": {"-123.456"},
				"float64": {"-789.012"},
			},
			want: TestStruct{
				Float32: -123.456,
				Float64: -789.012,
			},
			wantErr: false,
		},
		{
			name: "integer as float",
			values: url.Values{
				"float32": {"42"},
				"float64": {"100"},
			},
			want: TestStruct{
				Float32: 42.0,
				Float64: 100.0,
			},
			wantErr: false,
		},
		{
			name: "invalid float",
			values: url.Values{
				"float32": {"not_a_float"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestStruct
			dec := pform.NewDecoder(tt.values)
			err := dec.Decode(&got)

			if tt.wantErr {
				assert.Error(t, err)
				var parseErr pform.ParseError
				assert.True(t, errors.As(err, &parseErr))
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.want.Float32, got.Float32, 0.0001)
				assert.InDelta(t, tt.want.Float64, got.Float64, 0.0000001)
			}
		})
	}
}

// Test complex types
func TestDecoder_ComplexTypes(t *testing.T) {
	type TestStruct struct {
		Complex64  complex64  `form:"complex64"`
		Complex128 complex128 `form:"complex128"`
	}

	tests := []struct {
		name    string
		values  url.Values
		want    TestStruct
		wantErr bool
	}{
		{
			name: "valid complex numbers",
			values: url.Values{
				"complex64":  {"(1+2i)"},
				"complex128": {"(3.14+2.718i)"},
			},
			want: TestStruct{
				Complex64:  complex(1, 2),
				Complex128: complex(3.14, 2.718),
			},
			wantErr: false,
		},
		{
			name: "real only",
			values: url.Values{
				"complex64":  {"5"},
				"complex128": {"10.5"},
			},
			want: TestStruct{
				Complex64:  complex(5, 0),
				Complex128: complex(10.5, 0),
			},
			wantErr: false,
		},
		{
			name: "imaginary only",
			values: url.Values{
				"complex64":  {"3i"},
				"complex128": {"7.5i"},
			},
			want: TestStruct{
				Complex64:  complex(0, 3),
				Complex128: complex(0, 7.5),
			},
			wantErr: false,
		},
		{
			name: "invalid complex",
			values: url.Values{
				"complex64": {"not_complex"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestStruct
			dec := pform.NewDecoder(tt.values)
			err := dec.Decode(&got)

			if tt.wantErr {
				assert.Error(t, err)
				var parseErr pform.ParseError
				assert.True(t, errors.As(err, &parseErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// Test boolean parsing
func TestDecoder_BoolTypes(t *testing.T) {
	type TestStruct struct {
		Bool1 bool `form:"bool1"`
		Bool2 bool `form:"bool2"`
	}

	tests := []struct {
		name    string
		values  url.Values
		want    TestStruct
		wantErr bool
	}{
		{
			name: "true string",
			values: url.Values{
				"bool1": {"true"},
				"bool2": {"TRUE"},
			},
			want: TestStruct{
				Bool1: true,
				Bool2: true,
			},
			wantErr: false,
		},
		{
			name: "false string",
			values: url.Values{
				"bool1": {"false"},
				"bool2": {"FALSE"},
			},
			want: TestStruct{
				Bool1: false,
				Bool2: false,
			},
			wantErr: false,
		},
		{
			name: "numeric true",
			values: url.Values{
				"bool1": {"1"},
				"bool2": {"t"},
			},
			want: TestStruct{
				Bool1: true,
				Bool2: true,
			},
			wantErr: false,
		},
		{
			name: "numeric false",
			values: url.Values{
				"bool1": {"0"},
				"bool2": {"f"},
			},
			want: TestStruct{
				Bool1: false,
				Bool2: false,
			},
			wantErr: false,
		},
		{
			name: "invalid bool",
			values: url.Values{
				"bool1": {"not_a_bool"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestStruct
			dec := pform.NewDecoder(tt.values)
			err := dec.Decode(&got)

			if tt.wantErr {
				assert.Error(t, err)
				var parseErr pform.ParseError
				assert.True(t, errors.As(err, &parseErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// Test string fields
func TestDecoder_StringTypes(t *testing.T) {
	type TestStruct struct {
		Name    string `form:"name"`
		Email   string `form:"email"`
		Special string `form:"special"`
	}

	tests := []struct {
		name   string
		values url.Values
		want   TestStruct
	}{
		{
			name: "simple strings",
			values: url.Values{
				"name":  {"John Doe"},
				"email": {"john@example.com"},
			},
			want: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
		{
			name: "empty string",
			values: url.Values{
				"name":  {""},
				"email": {"test@example.com"},
			},
			want: TestStruct{
				Name:  "",
				Email: "test@example.com",
			},
		},
		{
			name: "special characters",
			values: url.Values{
				"name":    {"Test User"},
				"special": {"!@#$%^&*()_+-=[]{}|;':\",./<>?"},
			},
			want: TestStruct{
				Name:    "Test User",
				Special: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			},
		},
		{
			name: "unicode characters",
			values: url.Values{
				"name": {"日本語"},
			},
			want: TestStruct{
				Name: "日本語",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestStruct
			dec := pform.NewDecoder(tt.values)
			err := dec.Decode(&got)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Test struct tag options
func TestDecoder_StructTagOptions(t *testing.T) {
	t.Run("omitempty skips empty values", func(t *testing.T) {
		type TestStruct struct {
			Required string `form:"required"`
			Optional string `form:"optional,omitempty"`
		}

		values := url.Values{
			"required": {"value"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "value", got.Required)
		assert.Equal(t, "", got.Optional)
	})

	t.Run("required returns error when missing", func(t *testing.T) {
		type TestStruct struct {
			Name string `form:"name,required"`
		}

		values := url.Values{}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.Error(t, err)
		var reqErr pform.RequiredFieldError
		assert.True(t, errors.As(err, &reqErr))
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("required with value succeeds", func(t *testing.T) {
		type TestStruct struct {
			Name string `form:"name,required"`
		}

		values := url.Values{
			"name": {"John"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "John", got.Name)
	})

	t.Run("dash ignores field", func(t *testing.T) {
		type TestStruct struct {
			Name    string `form:"name"`
			Ignored string `form:"-"`
		}

		values := url.Values{
			"name":    {"John"},
			"Ignored": {"should be ignored"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "John", got.Name)
		assert.Equal(t, "", got.Ignored)
	})

	t.Run("empty tag ignores field", func(t *testing.T) {
		type TestStruct struct {
			Name      string `form:"name"`
			NoFormTag string
		}

		values := url.Values{
			"name":      {"John"},
			"NoFormTag": {"should be ignored"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "John", got.Name)
		assert.Equal(t, "", got.NoFormTag)
	})
}

// Test custom Unmarshaller interface
func TestDecoder_CustomUnmarshaller(t *testing.T) {
	t.Run("pointer receiver implementation", func(t *testing.T) {
		type TestStruct struct {
			Custom CustomType `form:"custom"`
		}

		values := url.Values{
			"custom": {"test data"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "custom:test data", got.Custom.Value)
	})

	t.Run("pointer field with Unmarshaller", func(t *testing.T) {
		type TestStruct struct {
			Custom *CustomType `form:"custom"`
		}

		values := url.Values{
			"custom": {"pointer data"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.NotNil(t, got.Custom)
		assert.Equal(t, "custom:pointer data", got.Custom.Value)
	})

	t.Run("Unmarshaller returning error", func(t *testing.T) {
		type TestStruct struct {
			Custom CustomErrorType `form:"custom"`
		}

		values := url.Values{
			"custom": {"error"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.Error(t, err)
		var parseErr pform.ParseError
		assert.True(t, errors.As(err, &parseErr))
		assert.Contains(t, err.Error(), "custom error")
	})

	t.Run("Unmarshaller with valid data", func(t *testing.T) {
		type TestStruct struct {
			Custom CustomErrorType `form:"custom"`
		}

		values := url.Values{
			"custom": {"valid"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "valid", got.Custom.Value)
	})
}

// Test map decoding
func TestDecoder_MapDecoding(t *testing.T) {
	t.Run("map[string]string", func(t *testing.T) {
		values := url.Values{
			"key1": {"value1"},
			"key2": {"value2"},
			"key3": {"value3"},
		}

		var got map[string]string
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		}, got)
	})

	t.Run("map[string]int", func(t *testing.T) {
		values := url.Values{
			"count1": {"10"},
			"count2": {"20"},
			"count3": {"30"},
		}

		var got map[string]int
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, map[string]int{
			"count1": 10,
			"count2": 20,
			"count3": 30,
		}, got)
	})

	t.Run("nil map initialization", func(t *testing.T) {
		values := url.Values{
			"key": {"value"},
		}

		var got map[string]string
		assert.Nil(t, got)
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "value", got["key"])
	})

	t.Run("map with invalid value type", func(t *testing.T) {
		values := url.Values{
			"key": {"not_a_number"},
		}

		var got map[string]int
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key")
	})

	t.Run("map with non-string key error", func(t *testing.T) {
		values := url.Values{
			"key": {"value"},
		}

		var got map[int]string
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "map key must be string")
	})

	t.Run("empty map", func(t *testing.T) {
		values := url.Values{}

		var got map[string]string
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got)
	})

	t.Run("map with multiple values uses first", func(t *testing.T) {
		values := url.Values{
			"key": {"first", "second", "third"},
		}

		var got map[string]string
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "first", got["key"])
	})
}

// Test HTTP request decoding (NewFormUrlDecoder)
func TestNewFormUrlDecoder(t *testing.T) {
	t.Run("valid form-urlencoded request", func(t *testing.T) {
		body := "name=John+Doe&age=30&email=john%40example.com"
		req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		decoder, err := pform.NewFormUrlDecoder(req)
		assert.NoError(t, err)
		assert.NotNil(t, decoder)

		type TestStruct struct {
			Name  string `form:"name"`
			Age   int    `form:"age"`
			Email string `form:"email"`
		}

		var got TestStruct
		err = decoder.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "John Doe", got.Name)
		assert.Equal(t, 30, got.Age)
		assert.Equal(t, "john@example.com", got.Email)
	})

	t.Run("Content-Type with charset parameter", func(t *testing.T) {
		body := "name=Test"
		req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

		decoder, err := pform.NewFormUrlDecoder(req)
		assert.NoError(t, err)
		assert.NotNil(t, decoder)

		type TestStruct struct {
			Name string `form:"name"`
		}

		var got TestStruct
		err = decoder.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "Test", got.Name)
	})

	t.Run("empty Content-Type error", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(""))
		assert.NoError(t, err)

		decoder, err := pform.NewFormUrlDecoder(req)
		assert.Error(t, err)
		assert.Nil(t, decoder)
		assert.Contains(t, err.Error(), "empty Content-Type")
	})

	t.Run("invalid Content-Type error", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(""))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "invalid/content/type/format")

		decoder, err := pform.NewFormUrlDecoder(req)
		assert.Error(t, err)
		assert.Nil(t, decoder)
		assert.Contains(t, err.Error(), "invalid Content-Type")
	})

	t.Run("wrong media type error", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(""))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		decoder, err := pform.NewFormUrlDecoder(req)
		assert.Error(t, err)
		assert.Nil(t, decoder)
		assert.Contains(t, err.Error(), "invalid Content-Type")
		assert.Contains(t, err.Error(), "application/x-www-form-urlencoded")
	})

	t.Run("mixed case Content-Type", func(t *testing.T) {
		body := "name=Test"
		req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "Application/X-WWW-FORM-URLENCODED")

		decoder, err := pform.NewFormUrlDecoder(req)
		assert.NoError(t, err)
		assert.NotNil(t, decoder)
	})

	t.Run("Content-Type with spaces", func(t *testing.T) {
		body := "name=Test"
		req, err := http.NewRequest("POST", "/test", bytes.NewBufferString(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "  application/x-www-form-urlencoded  ")

		decoder, err := pform.NewFormUrlDecoder(req)
		assert.NoError(t, err)
		assert.NotNil(t, decoder)
	})
}

// Test error conditions
func TestDecoder_ErrorConditions(t *testing.T) {
	t.Run("non-pointer destination", func(t *testing.T) {
		type TestStruct struct {
			Name string `form:"name"`
		}

		values := url.Values{"name": {"John"}}
		var obj TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(obj) // Not a pointer
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a pointer")
	})

	t.Run("invalid destination type", func(t *testing.T) {
		values := url.Values{"name": {"John"}}
		var obj int
		dec := pform.NewDecoder(values)
		err := dec.Decode(&obj)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pointer to struct, or map")
	})

	t.Run("type conversion error preserves context", func(t *testing.T) {
		type TestStruct struct {
			Age int `form:"age"`
		}

		values := url.Values{"age": {"not_a_number"}}
		var obj TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&obj)
		assert.Error(t, err)

		var parseErr pform.ParseError
		assert.True(t, errors.As(err, &parseErr))
		assert.Equal(t, "age", parseErr.Field)
		assert.Equal(t, "not_a_number", parseErr.Value)
	})
}

// Test edge cases
func TestDecoder_EdgeCases(t *testing.T) {
	t.Run("unexported struct fields are skipped", func(t *testing.T) {
		type TestStruct struct {
			Name       string `form:"name"`
			unexported string `form:"unexported"`
		}

		values := url.Values{
			"name":       {"John"},
			"unexported": {"should be ignored"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "John", got.Name)
		assert.Equal(t, "", got.unexported)
	})

	t.Run("multiple values for same key uses first", func(t *testing.T) {
		type TestStruct struct {
			Name string `form:"name"`
		}

		values := url.Values{
			"name": {"first", "second", "third"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "first", got.Name)
	})

	t.Run("empty form values", func(t *testing.T) {
		type TestStruct struct {
			Name string `form:"name"`
			Age  int    `form:"age"`
		}

		values := url.Values{}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, "", got.Name)
		assert.Equal(t, 0, got.Age)
	})

	t.Run("empty string value for integer field", func(t *testing.T) {
		type TestStruct struct {
			Age int `form:"age"`
		}

		values := url.Values{
			"age": {""},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.Error(t, err)
		var parseErr pform.ParseError
		assert.True(t, errors.As(err, &parseErr))
	})

	t.Run("mixed valid and invalid fields", func(t *testing.T) {
		type TestStruct struct {
			Name  string `form:"name"`
			Age   int    `form:"age"`
			Email string `form:"email"`
		}

		values := url.Values{
			"name":  {"John"},
			"age":   {"not_a_number"},
			"email": {"john@example.com"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.Error(t, err)
		// Name should have been set before the error occurred
		assert.Equal(t, "John", got.Name)
	})

	t.Run("zero values are valid", func(t *testing.T) {
		type TestStruct struct {
			Count  int     `form:"count"`
			Amount float64 `form:"amount"`
			Active bool    `form:"active"`
		}

		values := url.Values{
			"count":  {"0"},
			"amount": {"0.0"},
			"active": {"false"},
		}

		var got TestStruct
		dec := pform.NewDecoder(values)
		err := dec.Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, 0, got.Count)
		assert.Equal(t, 0.0, got.Amount)
		assert.Equal(t, false, got.Active)
	})
}

// Test the original basic test case
func TestDecoder_Decode(t *testing.T) {
	type TestStruct struct {
		Name  string     `form:"name"`
		Age   int        `form:"age"`
		Email string     `form:"email,omitempty"`
		Data  CustomType `form:"data"`
	}

	values := url.Values{
		"name": {"John Doe"},
		"age":  {"30"},
		"data": {"crazy data"},
	}

	var obj TestStruct
	dec := pform.NewDecoder(values)
	err := dec.Decode(&obj)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", obj.Name)
	assert.Equal(t, 30, obj.Age)
	assert.Equal(t, "", obj.Email) // Omitempty should not set a value
	assert.Equal(t, "custom:crazy data", obj.Data.Value)
}

// Benchmark tests
func BenchmarkDecoder_Struct(b *testing.B) {
	type TestStruct struct {
		Name  string  `form:"name"`
		Age   int     `form:"age"`
		Email string  `form:"email"`
		Score float64 `form:"score"`
	}

	values := url.Values{
		"name":  {"John Doe"},
		"age":   {"30"},
		"email": {"john@example.com"},
		"score": {"95.5"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var obj TestStruct
		dec := pform.NewDecoder(values)
		_ = dec.Decode(&obj)
	}
}

func BenchmarkDecoder_Map(b *testing.B) {
	values := url.Values{
		"key1": {"value1"},
		"key2": {"value2"},
		"key3": {"value3"},
		"key4": {"value4"},
		"key5": {"value5"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var obj map[string]string
		dec := pform.NewDecoder(values)
		_ = dec.Decode(&obj)
	}
}
