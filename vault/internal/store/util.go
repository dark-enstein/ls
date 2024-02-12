package store

const (
	ErrTokenTypeNotString = "token of type string required"
)

// InterfaceIsString checks if an interface underlying type is a string
func InterfaceIsString(a any) (bool, string) {
	switch v := a.(type) {
	case string:
		return true, v
	default:
		return false, ""
	}
}
