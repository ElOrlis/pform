package pf

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	tag       = "form"
	omitempty = "omitempty"
	required  = "required"
	omit      = "-"
	empty     = ""
)

func NewDecoder(form url.Values) Decoder {
	return Decoder{values: form}
}

type Decoder struct {
	values url.Values
}

func (d Decoder) Decode(dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Pointer {
		return errors.New("destination is not a pointer") // Create custom error
	}
	v = v.Elem()
	switch v.Kind() {
	case reflect.Struct:
		return d.decodeStruct(v)
	case reflect.Map:
		return nil
	default:
		return nil
	}
}

type Unmarshaller interface {
	UnmarshalForm(string) error
}

func setValueInt(field reflect.Value, v string) error {
	switch field.Kind() {
	case reflect.Int:
		d, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			return err
		}
		field.SetInt(d)
	case reflect.Int8:
		d, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return err
		}
		field.SetInt(d)
	case reflect.Int16:
		d, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return err
		}
		field.SetInt(d)
	case reflect.Int32:
		d, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(d)
	case reflect.Int64:
		d, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(d)
	default:
		return errors.New("field is not a valid integer type")
	}
	return nil
}

func setValueUint(field reflect.Value, v string) error {
	switch field.Kind() {
	case reflect.Uint:
		d, err := strconv.ParseUint(v, 10, 0)
		if err != nil {
			return err
		}
		field.SetUint(d)
	case reflect.Uint8:
		d, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return err
		}
		field.SetUint(d)
	case reflect.Uint16:
		d, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return err
		}
		field.SetUint(d)
	case reflect.Uint32:
		d, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return err
		}
		field.SetUint(d)
	case reflect.Uint64:
		d, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(d)
	default:
		return errors.New("field is not a valid unsigned integer type")
	}
	return nil
}

func setValueFloat(field reflect.Value, v string) error {
	switch field.Kind() {
	case reflect.Float32:
		d, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return err
		}
		field.SetFloat(d)
	case reflect.Float64:
		d, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		field.SetFloat(d)
	default:
		return errors.New("field is not a valid float type")
	}
	return nil
}

func setValueComplex(field reflect.Value, v string) error {
	switch field.Kind() {
	case reflect.Complex64:
		d, err := strconv.ParseComplex(v, 64)
		if err != nil {
			return err
		}
		field.SetComplex(d)
	case reflect.Complex128:
		d, err := strconv.ParseComplex(v, 128)
		if err != nil {
			return err
		}
		field.SetComplex(d)
	default:
		return errors.New("field is not a valid complex type")
	}
	return nil
}

func setValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setValueUint(field, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setValueInt(field, value)
	case reflect.Float32, reflect.Float64:
		return setValueFloat(field, value)
	case reflect.Complex64, reflect.Complex128:
		return setValueComplex(field, value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(v)
	default:
		ptr := reflect.New(field.Type())
		if !ptr.Type().Implements(reflect.TypeOf((*Unmarshaller)(nil)).Elem()) {
			return errors.New("field does not implement form Unmarshaller")
		}
		f, ok := ptr.Interface().(Unmarshaller)
		if !ok {
			return errors.New("field does not implement form Unmarshaller")
		}
		err := f.UnmarshalForm(value)
		if err != nil {
			return err
		}
		field.Set(ptr.Elem())
	}
	return nil
}

func (d Decoder) decodeStruct(v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}
		tagValue := v.Type().Field(i).Tag.Get(tag)
		if tagValue == omit || tagValue == empty {
			continue
		}
		schema := strings.Split(tagValue, ",")
		value := d.values.Get(schema[0])
		if len(value) == 0 {
			if len(schema) > 1 && schema[1] == omitempty {
				continue
			} else if len(schema) > 1 && schema[1] == required {
				return NewRequiredFieldError(schema[0])
			}
		}

		err := setValue(field, value)
		if err != nil {
			return err
		}
	}
	return nil
}
