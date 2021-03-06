package request

import (
	"net/http"

	"github.com/tidepool-org/platform/errors"
)

const (
	ErrorCodeInternalServerError = "internal-server-error"
	ErrorCodeUnexpectedResponse  = "unexpected-response"
	ErrorCodeTooManyRequests     = "too-many-requests"
	ErrorCodeBadRequest          = "bad-request"
	ErrorCodeUnauthenticated     = "unauthenticated"
	ErrorCodeUnauthorized        = "unauthorized"
	ErrorCodeResourceNotFound    = "resource-not-found"
	ErrorCodeHeaderMissing       = "header-missing"
	ErrorCodeHeaderInvalid       = "header-invalid"
	ErrorCodeParameterMissing    = "parameter-missing"
	ErrorCodeParameterInvalid    = "parameter-invalid"
	ErrorCodeJSONMalformed       = "json-malformed"
)

func ErrorInternalServerError(err error) error {
	return errors.WrapPrepared(err, ErrorCodeInternalServerError, "internal server error", "internal server error")
}

func ErrorUnexpectedResponse(res *http.Response, req *http.Request) error {
	return errors.Preparedf(ErrorCodeUnexpectedResponse, "unexpected response", "unexpected response status code %d from %s %q", res.StatusCode, req.Method, req.URL.String())
}

func ErrorTooManyRequests() error {
	return errors.Prepared(ErrorCodeTooManyRequests, "too many requests", "too many requests")
}

func ErrorBadRequest() error {
	return errors.Prepared(ErrorCodeBadRequest, "bad request", "bad request")
}

func ErrorUnauthenticated() error {
	return errors.Prepared(ErrorCodeUnauthenticated, "authentication token is invalid", "authentication token is invalid")
}

func ErrorUnauthorized() error {
	return errors.Prepared(ErrorCodeUnauthorized, "authentication token is not authorized for requested action", "authentication token is not authorized for requested action")
}

func ErrorResourceNotFound() error {
	return errors.Prepared(ErrorCodeResourceNotFound, "resource not found", "resource not found")
}

func ErrorResourceNotFoundWithID(id string) error {
	return errors.Preparedf(ErrorCodeResourceNotFound, "resource not found", "resource with id %q not found", id)
}

func ErrorResourceNotFoundWithIDAndRevision(id string, revision int) error {
	return errors.Preparedf(ErrorCodeResourceNotFound, "resource not found", "revision %d of resource with id %q not found", revision, id)
}

func ErrorHeaderMissing(key string) error {
	return errors.Preparedf(ErrorCodeHeaderMissing, "header is missing", "header %q is missing", key)
}

func ErrorHeaderInvalid(key string) error {
	return errors.Preparedf(ErrorCodeHeaderInvalid, "header is invalid", "header %q is invalid", key)
}

func ErrorParameterMissing(key string) error {
	return errors.Preparedf(ErrorCodeParameterMissing, "parameter is missing", "parameter %q is missing", key)
}

func ErrorParameterInvalid(key string) error {
	return errors.Preparedf(ErrorCodeParameterInvalid, "parameter is invalid", "parameter %q is invalid", key)
}

func ErrorJSONMalformed() error {
	return errors.Prepared(ErrorCodeJSONMalformed, "json is malformed", "json is malformed")
}

func StatusCodeForError(err error) int {
	if err != nil {
		switch errors.Code(err) {
		case ErrorCodeTooManyRequests:
			return http.StatusTooManyRequests
		case ErrorCodeBadRequest:
			return http.StatusBadRequest
		case ErrorCodeUnauthenticated:
			return http.StatusUnauthorized
		case ErrorCodeUnauthorized:
			return http.StatusForbidden
		case ErrorCodeResourceNotFound:
			return http.StatusNotFound
		}
	}
	return http.StatusInternalServerError
}

func IsErrorInternalServerError(err error) bool {
	return errors.Code(err) == ErrorCodeInternalServerError
}

func IsErrorUnauthenticated(err error) bool {
	return errors.Code(err) == ErrorCodeUnauthenticated
}

func IsErrorUnauthorized(err error) bool {
	return errors.Code(err) == ErrorCodeUnauthorized
}

func IsErrorResourceNotFound(err error) bool {
	return errors.Code(err) == ErrorCodeResourceNotFound
}
