package errors

const (
	None Type = 0

	BadRequest Type = 4000 + iota
	Unauthorized
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
		return "unknown"
	}
}
