package errors

type ErrorString string

/*
ErrorStrings are errors
*/
func (e ErrorString) Error() string {
	return string(e)
}

const (
	NoMatchingEntityInstanceFound ErrorString = "no matching entity instance found"
	EntityTypeHasNoDescendants    ErrorString = "entity type has no descendants"
)
