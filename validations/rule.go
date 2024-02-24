package validations

// Rule represents a validation rule.
// This is compatible with github.com/go-ozzo/ozzo-validation/v4.Rule.
type Rule interface {
	// Validate validates a value and returns a value if validation fails.
	Validate(value interface{}) error
}
