// Package inputconv converts nullable generated API inputs into values accepted
// by ent's SetNillable methods.
package inputconv

// Nullable is the common interface implemented by ogen's generated Nil* types.
type Nullable[T any] interface {
	Get() (T, bool)
}

// Ptr returns a pointer to the wrapped value, or nil when the input is null.
func Ptr[T any](value Nullable[T]) *T {
	v, ok := value.Get()
	if !ok {
		return nil
	}
	return &v
}

// StringPtr converts a nullable string-backed value, including generated enums,
// into the plain string pointer expected by ent string setters.
func StringPtr[T ~string](value Nullable[T]) *string {
	v, ok := value.Get()
	if !ok {
		return nil
	}
	s := string(v)
	return &s
}

// IntPtr converts a nullable contract int32 into the native int pointer
// expected by ent integer setters.
func IntPtr[T ~int32](value Nullable[T]) *int {
	v, ok := value.Get()
	if !ok {
		return nil
	}
	n := int(v)
	return &n
}
