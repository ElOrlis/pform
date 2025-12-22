package pform

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	tag                       = "form"
	omitempty                 = "omitempty"
	required                  = "required"
	omit                      = "-"
	empty                     = ""
	formUrlEncodedContentType = "application/x-www-form-urlencoded"
)

type Unmarshaller interface {
	UnmarshalValue(string) error
}

type XFormUrlEncoded interface {
	Decode(dest any) error
}

func NewFormUrlDecoder(r *http.Request) (XFormUrlEncoded, error) {
	contentType := r.Header.Get("Content-Type")
	contentType = strings.TrimSpace(contentType)
	contentType = strings.ToLower(contentType)

	if len(contentType) == 0 {
		return nil, errors.New("empty Content-Type in HTTP request")
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Type: %w", err)
	}

	if mediaType != formUrlEncodedContentType {
		return nil, fmt.Errorf(
			"invalid Content-Type %q, expected application/x-www-form-urlencoded",
			mediaType,
		)
	}

	err = r.ParseForm()
	if err != nil {
		return nil, err
	}

	return NewDecoder(r.PostForm), nil
}

func NewDecoder(form url.Values) Decoder {
	return Decoder{values: form}
}

type Decoder struct {
	values url.Values
}

func (d Decoder) Decode(dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Pointer {
		return errors.New("interface is not a pointer")
	}
	v = v.Elem()
	switch v.Kind() {
	case reflect.Struct:
		return d.decodeStruct(v)
	case reflect.Map:
		return d.decodeMap(v)
	default:
		return errors.New("interface should be a pointer to struct, or map")
	}
}

func (d Decoder) decodeStruct(v reflect.Value) error {
	for i := range v.NumField() {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		tagValue := v.Type().Field(i).Tag.Get(tag)
		if tagValue == omit || tagValue == empty {
			continue
		}

		schema := strings.Split(tagValue, ",")
		values, exists := d.values[schema[0]]

		if !exists || len(values) == 0 {
			if len(schema) > 1 && schema[1] == required {
				return newRequiredFieldError(schema[0])
			}
			continue
		}

		value := values[0]

		if len(value) == 0 {
			if len(schema) > 1 && schema[1] == omitempty {
				continue
			} else if len(schema) > 1 && schema[1] == required {
				return newRequiredFieldError(schema[0])
			}
		}

		err := setValue(field, value, schema[0])
		if err != nil {
			return err
		}
	}
	return nil
}

func (d Decoder) decodeMap(v reflect.Value) error {
	if v.Type().Key().Kind() != reflect.String {
		return errors.New("map key must be string type")
	}

	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}

	elemType := v.Type().Elem()
	for key, values := range d.values {
		if len(values) == 0 {
			continue
		}

		elem := reflect.New(elemType).Elem()
		err := setValue(elem, values[0], key)
		if err != nil {
			return fmt.Errorf("error setting map value for key %q: %w", key, err)
		}

		v.SetMapIndex(reflect.ValueOf(key), elem)
	}
	return nil
}
