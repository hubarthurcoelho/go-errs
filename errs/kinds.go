package errs

import "net/http"

type kind uint8

// Kinds of errors.
//
// Do not reorder this list or remove any
// items since that will change their values.
// New items must be added only to the end.
const (
	UnauthorizedCredential kind = iota + 1 // credential not active
	InvalidCredential                      // credential not found
	RepositoryError                        // error in our repository
	SupplierError                          // error in some of our suppliers
	ValidationError                        // data validation failed
	JSONError                              // marshal/unmarshal/readAll operations failed
	InvalidInput                           // error in internal packages, like utils
	NotFound                               // some record was not found
	Unprocessable                          // cannot proceed with some request
	CacheMalfunction                       // something cache-related failed
)

func (k kind) String() string {
	switch k {
	case UnauthorizedCredential:
		return "UNAUTHORIZED_CREDENTIAL"
	case InvalidCredential:
		return "INVALID_CREDENTIAL"
	case RepositoryError:
		return "REPOSITORY_ERROR"
	case SupplierError:
		return "SUPPLIER_ERROR"
	case ValidationError:
		return "VALIDATION_ERROR"
	case JSONError:
		return "JSON_ERROR"
	case InvalidInput:
		return "INVALID_INPUT"
	case NotFound:
		return "NOT_FOUND"
	case CacheMalfunction:
		return "CACHE_ERROR"
	default:
		return "UNEXPECTED_ERROR"
	}
}

func (k kind) HttpStatus() int {
	switch k {
	case InvalidCredential:
		return http.StatusBadRequest
	case UnauthorizedCredential:
		return http.StatusUnauthorized
	case RepositoryError:
		return http.StatusInternalServerError
	case SupplierError:
		return http.StatusBadGateway
	case ValidationError:
		return http.StatusBadRequest
	case JSONError:
		return http.StatusInternalServerError
	case InvalidInput:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case CacheMalfunction:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
