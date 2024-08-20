package models

type Model interface {
	// String containing field names to use when querying a DB
	FieldNames() string
	// Slice of current values of fields in Model
	Values() []any

	Fields
}

type Fields interface {
	// Slice of references to fields
	Fields() []any
}

func ValuesFromFields(f Fields) []any {
	fields := f.Fields()
	values := make([]any, len(fields))
	for i, f := range fields {
		// Cast f into *any and dereference
		v := *(f.(*any))
		values[i] = v
	}
	return values
}
