// Package types defines shared types and response structures for the application.
package types

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitresponse"`
}
