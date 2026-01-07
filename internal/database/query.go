package database

import (
	"fmt"
	"reflect"
	"strings"
)

var allowedOps = map[string]bool{
	"=": true, "!=": true, "<>": true,
	"<": true, "<=": true, ">": true, ">=": true,
	"like": true, "ilike": true, "in": true,
}

func (db *DB) BuildQueryArgs(filter any) ([]string, []interface{}, error) {
	if filter == nil {
		return nil, nil, nil
	}

	v := reflect.ValueOf(filter)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, nil, nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("filter must be a struct or *struct, got %s", v.Kind())
	}

	t := v.Type()
	whereClauses := make([]string, 0, t.NumField())
	args := make([]any, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := v.Field(i)

		// skip unexported fields
		if sf.PkgPath != "" {
			continue
		}

		col := sf.Tag.Get("db")
		if col == "" || col == "-" {
			continue
		}
		if !isSafeIdentifier(col) {
			return nil, nil, fmt.Errorf("unsafe db tag column: %q", col)
		}

		op := sf.Tag.Get("op")
		if op == "" {
			op = "="
		}
		opLower := strings.ToLower(op)

		// optional-by-pointer
		if fv.Kind() == reflect.Pointer {
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
		}

		// IN support
		if opLower == "in" {
			if fv.Kind() != reflect.Slice && fv.Kind() != reflect.Array {
				return nil, nil, fmt.Errorf("%s has op=in but is %s", sf.Name, fv.Kind())
			}
			if fv.Len() == 0 {
				continue
			}

			phs := make([]string, 0, fv.Len())
			for j := 0; j < fv.Len(); j++ {
				phs = append(phs, "?")
				args = append(args, fv.Index(j).Interface())
			}
			whereClauses = append(
				whereClauses,
				fmt.Sprintf("%s IN (%s)", col, strings.Join(phs, ", ")),
			)
			continue
		}

		// scalar predicate
		whereClauses = append(
			whereClauses,
			fmt.Sprintf("%s %s ?", col, op),
		)
		args = append(args, fv.Interface())
	}

	return whereClauses, args, nil
}

func isSafeIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '_' || r == '.':
		default:
			return false
		}
	}
	return true
}
