package middleware

import (
	"net/http"
)

func responseCodeToString(responseCode int) string {
	return http.StatusText(responseCode)
}
