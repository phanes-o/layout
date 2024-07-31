package errors

// common Error type
const (
	None Type = 0

	Unauthorized Type = 1000

	ErrParamsParse Type = iota + 1000
)

const (
	BadRequest Type = iota + 2000
	Forbidden
	NotFound
	InternalServerError
	NotImplemented
	ServiceUnavailable
)

func (t Type) String() string {
	switch t {
	case None:
		return "none"
	case BadRequest:
		return "bad request"
	case Unauthorized:
		return "unauthorized"
	case ErrParamsParse:
		return "ErrParamsParse"

	case Forbidden:
		return "forbidden"
	case NotFound:
		return "not found"
	case InternalServerError:
		return "internal server errors"
	case NotImplemented:
		return "not implemented"
	case ServiceUnavailable:
		return "service unavailable"
	default:
		return ""
	}
}
