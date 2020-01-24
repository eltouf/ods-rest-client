package odsrestclient

import "fmt"

//RateLimitError : Api Quotas https://help.opendatasoft.com/apis/ods-search-v1/#quotas
type RateLimitError struct {
	Limit     uint16
	Remaining uint16
	Reset     uint16
}

func (err *RateLimitError) Error() string {
	return fmt.Sprintf("%d of %d remaining requests. Next reset at %v", err.Remaining, err.Limit, err.Reset)
}
