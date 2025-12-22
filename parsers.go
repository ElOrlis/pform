package pform

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func setValueInt(field reflect.Value, v, n string) error {
	t := field.Type().Kind().String()
	switch field.Kind() {
	case reflect.Int:
		d, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetInt(d)
	case reflect.Int8:
		d, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetInt(d)
	case reflect.Int16:
		d, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetInt(d)
	case reflect.Int32:
		d, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetInt(d)
	case reflect.Int64:
		d, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetInt(d)
	default:
		return newParseError(n, t, v, errors.New("value is not an integer"))
	}
	return nil
}

func setValueUint(field reflect.Value, v, n string) error {
	t := field.Type().Kind().String()
	switch field.Kind() {
	case reflect.Uint:
		d, err := strconv.ParseUint(v, 10, 0)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetUint(d)
	case reflect.Uint8:
		d, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetUint(d)
	case reflect.Uint16:
		d, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetUint(d)
	case reflect.Uint32:
		d, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetUint(d)
	case reflect.Uint64:
		d, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetUint(d)
	default:
		return newParseError(n, t, v, errors.New("value is not an unsigned integer"))
	}
	return nil
}

func setValueFloat(field reflect.Value, v, n string) error {
	t := field.Type().Kind().String()
	switch field.Kind() {
	case reflect.Float32:
		d, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetFloat(d)
	case reflect.Float64:
		d, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetFloat(d)
	default:
		return newParseError(n, t, v, errors.New("value is not a float"))
	}
	return nil
}

func setValueComplex(field reflect.Value, v, n string) error {
	t := field.Type().Kind().String()
	switch field.Kind() {
	case reflect.Complex64:
		d, err := strconv.ParseComplex(v, 64)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetComplex(d)
	case reflect.Complex128:
		d, err := strconv.ParseComplex(v, 128)
		if err != nil {
			return newParseError(n, t, v, err)
		}
		field.SetComplex(d)
	default:
		return newParseError(n, t, v, errors.New("value is not a complex"))
	}
	return nil
}

func setValue(field reflect.Value, value, name string) error {
	t := field.Type().Kind().String()
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setValueUint(field, value, name)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setValueInt(field, value, name)
	case reflect.Float32, reflect.Float64:
		return setValueFloat(field, value, name)
	case reflect.Complex64, reflect.Complex128:
		return setValueComplex(field, value, name)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return newParseError(name, t, value, err)
		}
		field.SetBool(v)
		return nil
	default:
		if reflect.PointerTo(field.Type()).Implements(reflect.TypeFor[Unmarshaller]()) {
			ptr := reflect.New(field.Type())
			unmarshaller := ptr.Interface().(Unmarshaller)
			if err := unmarshaller.UnmarshalValue(value); err != nil {
				return newParseError(name, t, value, err)
			}
			field.Set(ptr.Elem())
			return nil
		}

		if field.Kind() == reflect.Pointer {
			elemType := field.Type().Elem()
			if reflect.PointerTo(elemType).Implements(reflect.TypeFor[Unmarshaller]()) {
				if field.IsNil() {
					field.Set(reflect.New(elemType))
				}
				unmarshaller := field.Interface().(Unmarshaller)
				if err := unmarshaller.UnmarshalValue(value); err != nil {
					return newParseError(name, t, value, err)
				}
				return nil
			}
		}
	}
	return newParseError(
		name,
		t,
		value,
		fmt.Errorf("unsupported field type: %s", t),
	)
}
