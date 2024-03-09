package jury

// Internal errors
const (
	ErrInternalNone = iota
	ErrInternalGeneric
)

// External errors
const (
	ErrSuccess = iota
	ErrGeneric
	ErrInternal
	ErrRequest
)

const (
	ErrUndefined = -1 // ideally should not be preferred. this is just for lazy allowance. will be deprecated with time
)
