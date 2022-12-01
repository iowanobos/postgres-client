package postgres

import (
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
)

const dbStructTag = "db"

type scanner struct {
	value                reflect.Value
	elem                 reflect.Type
	positionByColumnName map[string]int
}

func newScanner(dest any) (*scanner, error) {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("postgres: unsupported type: %s. should use a pointer of slice", value.String())
	}
	value = value.Elem()

	if value.Kind() == reflect.Slice {
		sliceValue := value.Type().Elem()

		if sliceValue.Kind() == reflect.Struct {
			// struct slice
			return &scanner{
				value:                value,
				elem:                 sliceValue,
				positionByColumnName: getPositionByColumnName(sliceValue),
			}, nil
		} else {
			// primitive slice
			return &scanner{
				value: value,
				elem:  sliceValue,
			}, nil
		}
	} else {
		if value.Kind() == reflect.Struct {
			// struct value
			return &scanner{
				value:                value,
				positionByColumnName: getPositionByColumnName(value.Type()),
			}, nil
		} else {
			// primitive value
			return &scanner{
				value: value,
			}, nil
		}
	}
}

// getPositionByColumnName parse field position by column name
func getPositionByColumnName(t reflect.Type) map[string]int {
	positionByColumnName := make(map[string]int, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		if columnName, ok := t.Field(i).Tag.Lookup(dbStructTag); ok {
			positionByColumnName[columnName] = i
		}
	}
	return positionByColumnName
}

func (s *scanner) scan(rows pgx.Rows) error {
	if s.elem != nil {
		if s.positionByColumnName != nil {
			// struct slice
			elem := newElem(s.elem)
			mapper, err := s.scanStructFirstRow(elem, rows)
			if err != nil {
				return err
			}
			appendSlice(s.value, elem)

			for rows.Next() {
				elem := newElem(s.elem)
				if err := scanStructByMapper(elem, rows, mapper); err != nil {
					return err
				}
				appendSlice(s.value, elem)
			}
		} else {
			// primitive slice
			for rows.Next() {
				elem := newElem(s.elem)
				if err := rows.Scan(getPointer(elem)); err != nil {
					return err
				}
				appendSlice(s.value, elem)
			}
		}
	} else {
		if s.positionByColumnName != nil {
			// struct value
			_, err := s.scanStructFirstRow(s.value, rows)
			return err
		} else {
			// primitive value
			if rows.Next() {
				if err := rows.Scan(getPointer(s.value)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *scanner) scanStructFirstRow(elem reflect.Value, rows pgx.Rows) (map[int]int, error) {
	mapper := make(map[int]int)
	if rows.Next() {
		for columnNumber, fieldDescription := range rows.FieldDescriptions() {
			if position, ok := s.positionByColumnName[fieldDescription.Name]; ok {
				mapper[position] = columnNumber
			}
		}

		if err := scanStructByMapper(elem, rows, mapper); err != nil {
			return nil, err
		}
	}
	return mapper, nil
}

func scanStructByMapper(elem reflect.Value, rows pgx.Rows, mapper map[int]int) error {
	values := make([]any, len(mapper))
	for i := 0; i < len(rows.FieldDescriptions()); i++ {
		if position, ok := mapper[i]; ok {
			values[position] = getPointer(elem.Field(i))
		}
	}

	if err := rows.Scan(values...); err != nil {
		return err
	}
	return nil
}

func newElem(t reflect.Type) reflect.Value {
	elem := reflect.New(t)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	return elem
}

func getPointer(v reflect.Value) any {
	if v.Kind() == reflect.Ptr {
		v.Set(reflect.New(v.Type().Elem()))
	} else {
		v = v.Addr()
	}
	return v.Interface()
}

func appendSlice(sliceValue reflect.Value, elem reflect.Value) {
	index := sliceValue.Len()
	if index >= sliceValue.Cap() {
		capacity := sliceValue.Cap() + sliceValue.Cap()/2
		if capacity < 4 {
			capacity = 4
		}
		slice := reflect.MakeSlice(sliceValue.Type(), sliceValue.Len(), capacity)
		reflect.Copy(slice, sliceValue)
		sliceValue.Set(slice)
	}
	if index >= sliceValue.Len() {
		sliceValue.SetLen(index + 1)
	}
	sliceValue.Index(index).Set(elem)
}
