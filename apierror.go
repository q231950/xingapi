/*
Package xingapi is a wrapper and connector to the XING API (https://dev.xing.com/docs/resources)
*/
package xingapi

// An APIError represents an error that occurred in the XING API error domain
type APIError struct {
	Message       string `json:"message"`
	ThrottledType string `json:"throttled"`
	BanTime       string `json:"ban_time"`
	ErrorName     string `json:"error_name"`
}

// String makes APIError conform to Stringer interface.
func (error APIError) String() string {
	return error.Message + ". Ban time: " + error.BanTime + " seconds"
}
