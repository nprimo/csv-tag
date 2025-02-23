package main

import (
	"fmt"
	"reflect"
)

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "csv: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "csv: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "csv: Unmarshal(nil " + e.Type.String() + ")"
}

// Unmarshal parse the content of a CSV file and store the result into the
// value pointed to by v
func Unmarshal(rows [][]string, v any) error {
	if reflect.TypeOf(v).Kind() != reflect.Pointer {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	if reflect.TypeOf(v).Elem().Kind() != reflect.Slice {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	header := rows[0]
	if err := checkValidHeader(header, v); err != nil {
		return err
	}
	d := newDecoder()
	if err := d.init(header, v); err != nil {
		return err
	}
	return d.unmarshall(rows, v)
}

// checkValidHeader check if the header has matching tags in the struct
// passed as v.
func checkValidHeader(header []string, v any) error {
	ptrEl := reflect.TypeOf(v).Elem()
	t := ptrEl.Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("csv")
		if i := index(tag, header); i == -1 {
			return fmt.Errorf("tag(%s) not present in header %+v",
				tag, header,
			)
		}
	}
	return nil
}

type decoder struct {
	mapper map[int]string
}

func newDecoder() decoder {
	return decoder{
		mapper: map[int]string{},
	}
}

func (d *decoder) init(header []string, v any) error {
	ptrEl := reflect.TypeOf(v).Elem()
	t := ptrEl.Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("csv")
		id := index(tag, header)
		d.mapper[id] = field.Name
	}
	return nil
}

// unmarshall store the data inside the array of struct v.
func (d *decoder) unmarshall(data [][]string, v any) error {
	sliceRv := reflect.MakeSlice(
		reflect.TypeOf(v).Elem(),
		len(data)-1,
		len(data)-1,
	)
	for i, row := range data[1:] {
		rv := sliceRv.Index(i)
		for rowIndex, fieldName := range d.mapper {
			f := rv.FieldByName(fieldName)
			switch f.Kind() {
			case reflect.String:
				f.SetString(row[rowIndex])
			default:
				// TODO: other cases
			}
		}
	}
	reflect.ValueOf(v).Elem().Set(sliceRv)
	return nil
}

// utils function to get the index of first occurrence of a string in an array
func index(target string, ss []string) int {
	for i, s := range ss {
		if s == target {
			return i
		}
	}
	return -1
}
