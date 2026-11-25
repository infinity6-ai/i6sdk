package tabledataz

import (
	"cmp"
	"fmt"
	"time"

	"github.com/infinity6-ai/gox/commonz/validation"
)

// Field validators for TableData.Validate implementations.
//
// Each one returns nil or an error formatted as "failed to validate <field>: <reason>",
// so a Validate method can be written as a single FirstError(...) call.
// Null* variants only validate when the value is present (Valid).

// FirstError returns the first non-nil error, or nil.
func FirstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func fieldErr(field string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("failed to validate %s: %w", field, err)
}

// Required checks that a string is not empty.
func Required(field, v string) error {
	return fieldErr(field, validation.StrNotEmpty(v, "%s is required", field))
}

// RequiredTime checks that a timestamp is set.
func RequiredTime(field string, v time.Time) error {
	return fieldErr(field, validation.False(v.IsZero(), "%s is required", field))
}

// Str checks that a string is not empty and has at most maxLen characters.
func Str(field, v string, maxLen int) error {
	return FirstError(Required(field, v), MaxLen(field, v, maxLen))
}

// MaxLen checks that a string has at most maxLen characters.
func MaxLen(field, v string, maxLen int) error {
	return fieldErr(field, validation.LessOrEqual(len(v), maxLen, "%s must have at most %d characters", field, maxLen))
}

// NullMaxLen checks the length of an optional string when present.
func NullMaxLen(field string, v *string, maxLen int) error {
	if v == nil {
		return nil
	}
	return MaxLen(field, *v, maxLen)
}

// OneOf checks that a string is one of the allowed choices.
func OneOf(field, v string, choices ...string) error {
	return fieldErr(field, validation.OneOf(choices, v, "%s must be one of: %v", field, choices))
}

// Between checks that min <= v <= max.
func Between[T cmp.Ordered](field string, v, min, max T) error {
	return fieldErr(field, FirstError(
		validation.GreaterOrEqual(v, min, "%s must be >= %v", field, min),
		validation.LessOrEqual(v, max, "%s must be <= %v", field, max),
	))
}

// NullIntBetween checks the range of an optional integer when present.
func NullIntBetween(field string, v *int64, min, max int64) error {
	if v == nil {
		return nil
	}
	return Between(field, *v, min, max)
}

// Flag checks that an integer flag is 0 or 1.
func Flag(field string, v int64) error {
	return Between(field, v, 0, 1)
}

// NullFlag checks that an optional integer flag is 0 or 1 when present.
func NullFlag(field string, v *int64) error {
	return NullIntBetween(field, v, 0, 1)
}
