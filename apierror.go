/*
xingapi is a wrapper and connector to the XING API (https://dev.xing.com/docs/resources)
*/
package xingapi

import "strconv"

// An APIError represents an error that occurred in the XING API error domain
type APIError struct {
	Message       string `json:"message"`
	ThrottledType string `json:"throttled"`
	BanTime       int    `json:"ban_time"`
	ErrorName     string `json:"error_name"`
}

func (error APIError) String() string {
	return error.Message + ". Ban time: " + strconv.Itoa(error.BanTime) + " seconds"
}
