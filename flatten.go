package flattenbson

import (
	"errors"
	"reflect"
	"strings"
)

// eg.
//
// type X struct {
// 	Quantity int `bson:"quantity"`
// }

// type Human struct {
// 	Eyes   *X `bson:"eyes,omitempty"`
// 	Tongue X  `bson:"tongue"`
// }

// Flatten(Human{Eyes: &X{Quantity: 2}, Tongue: X{}})
// result > map[string]any{"eyes.quantity" : 2, "tongue.quantity": 0}

// Recursively walks the struct and flattens it into a map, supported bson tags are omitempty, minsize, truncate, -
func Flatten(data any) (map[string]any, error) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("wrong type! want type struct in argument")
	}

	finalMap := make(map[string]any)

	for i := range v.NumField() {
		fv := v.Field(i)
		ft := t.Field(i)

		if !ft.IsExported() {
			continue
		}

		tag := ft.Tag.Get("bson")
		nameAndOptions := strings.Split(tag, ",")
		var name string = ft.Name
		if nameAndOptions[0] != "" {
			name = nameAndOptions[0]
		}

		if fv.Kind() == reflect.Pointer {
			fv = fv.Elem()
		}

		// TODO: process tags
		if fv.Kind() == reflect.Struct {
			if !fv.IsValid() {
				continue
			}
			m, err := Flatten(fv.Interface())
			if err != nil {
				return nil, err
			}

			err = addPrefix(name, m, finalMap)
			if err != nil {
				return nil, err
			}
			continue
		}

		if !fv.IsValid() {
			continue
		}
		finalMap[name] = fv.Interface()
	}

	return finalMap, nil
}

func addPrefix(preKey string, src map[string]any, dest map[string]any) error {
	for k, v := range src {
		dest[preKey+"."+k] = v
	}
	return nil
}
