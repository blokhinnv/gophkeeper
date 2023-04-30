package encrypt

import (
	"fmt"
	"reflect"
)

func EncryptStringFields(data any, key string) (any, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not a pointer: %v", v.Kind())
	}
	s := v.Elem()
	if s.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct")
	}
	for i := 0; i < s.NumField(); i++ {
		field, ok := s.Field(i).Interface().(string)
		if !ok {
			continue
		}
		encField, err := EncryptString(field, key)
		if err != nil {
			return nil, err
		}
		s.Field(i).SetString(encField)
	}
	return data, nil
}

func DecryptStringFields(data any, key string) (any, error) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("not a pointer: %v", v.Kind())
	}
	s := v.Elem()
	if s.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct: %v", s.Kind())
	}
	for i := 0; i < s.NumField(); i++ {
		field, ok := s.Field(i).Interface().(string)
		if !ok {
			continue
		}
		encField, err := DecryptString(field, key)
		if err != nil {
			return nil, err
		}
		s.Field(i).SetString(encField)
	}
	return data, nil
}
