package utils

// SetIfNotEmpty sets dst to src only if src is not an empty string
func SetIfNotEmpty(dst *string, src string) {
	if src != "" {
		*dst = src
	}
}

// SetIfNotEmptyStr is a generic version for string-like types (string, custom string types)
func SetIfNotEmptyStr[T ~string](dst *T, src T) {
	if src != "" {
		*dst = src
	}
}

// SetIfPositive sets dst to src only if src is greater than 0
func SetIfPositive(dst *int, src int) {
	if src > 0 {
		*dst = src
	}
}

// SetIfNotNil sets dst to the value of src only if src is not nil
func SetIfNotNil[T any](dst *T, src *T) {
	if src != nil {
		*dst = *src
	}
}

// Of returns a pointer to the given value
func Of[T any](v T) *T {
	return &v
}

// Deref returns the value pointed to by p, or the zero value if p is nil
func Deref[T any](p *T) T {
	if p != nil {
		return *p
	}
	var zero T
	return zero
}

// DerefOr returns the value pointed to by p, or the provided default value if p is nil
func DerefOr[T any](p *T, defaultVal T) T {
	if p != nil {
		return *p
	}
	return defaultVal
}
