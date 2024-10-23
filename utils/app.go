package utils

func init() {
	In.Array = InArray[any]
	Array.Diff = ArrayDiff[any]
	Array.Filter = ArrayFilter
	Array.Remove = ArrayRemove
	Array.Unique = ArrayUnique[any]
	Array.Empty = ArrayEmpty[any]
	Array.Reverse = ArrayReverse[any]
	Array.MapWithField = ArrayMapWithField
	Map.WithField = MapWithField[map[string]any]
	Map.WithoutField = MapWithoutField[map[string]any]
	Map.ToURL  = MapToURL
	Map.Keys   = MapKeys[map[string]any]
	Map.Values = MapValues[map[string]any]
	Map.Trim   = MapTrim[map[string]any]
}

var In struct {
	Array func(value any, array []any) (ok bool)
}

var Array struct {
	Diff         func(array1, array2 []any) (slice []any)
	Filter       func(array []string) (slice []string)
	Remove       func(array []string, args ...string) (slice []string)
	Unique       func(array []any) (slice []any)
	Empty        func(array []any) (slice []any)
	Reverse      func(array []any) (slice []any)
	MapWithField func(array []map[string]any, field any) (slice []any)
}

var Map struct {
	WithField    func(data map[string]any, field []string) (result map[string]any)
	WithoutField func(data map[string]any, field []string) (result map[string]any)
	ToURL        func(data map[string]any) (result string)
	Keys         func(data map[string]any) (result []string)
	Values       func(data map[string]any) (result []any)
	Trim	     func(data map[string]any) (result map[string]any)
}
