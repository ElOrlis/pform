package pform_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"pform"
)

type CustomType struct {
	Value string
}

func (c *CustomType) UnmarshalValue(v string) error {
	c.Value = "custom:" + v
	return nil
}

// Test Decoder with struct
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
