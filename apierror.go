// apierror.go
package xingapi

import "strconv"

type APIError struct {
	Message       string `json:"message"`
	ThrottledType string `json:"throttled"`
	BanTime       int    `json:"ban_time"`
	ErrorName     string `json:"error_name"`
}

func (error APIError) String() string {
	return error.Message + ". Ban time: " + strconv.Itoa(error.BanTime) + " seconds"
}
