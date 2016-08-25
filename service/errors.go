package service

/* CHECKLIST
 * [x] Uses interfaces as appropriate
 * [x] Private package variables use underscore prefix
 * [x] All parameters validated
 * [x] All errors handled
 * [x] Reviewed for concurrency safety
 * [x] Code complete
 * [x] Full test coverage
 */

import "net/http"

func ErrorInternalServerFailure() *Error {
	return &Error{
		Code:   "internal-server-failure",
		Status: http.StatusInternalServerError,
		Title:  "internal server failure",
		Detail: "Internal server failure",
	}
}

func ErrorJSONMalformed() *Error {
	return &Error{
		Code:   "json-malformed",
		Status: http.StatusBadRequest,
		Title:  "json is malformed",
		Detail: "JSON is malformed",
	}
}

func ErrorAuthenticationTokenMissing() *Error {
	return &Error{
		Code:   "authentication-token-missing",
		Status: http.StatusUnauthorized,
		Title:  "authentication token missing",
		Detail: "Authentication token missing",
	}
}

func ErrorUnauthenticated() *Error {
	return &Error{
		Code:   "unauthenticated",
		Status: http.StatusUnauthorized,
		Title:  "authentication token is invalid",
		Detail: "Authentication token is invalid",
	}
}

func ErrorUnauthorized() *Error {
	return &Error{
		Code:   "unauthorized",
		Status: http.StatusForbidden,
		Title:  "authentication token is not authorized for requested action",
		Detail: "Authentication token is not authorized for requested action",
	}
}
