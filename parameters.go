package odsrestclient

import "net/url"

// ODSParameters interface to convert parameters into url.Values type
type ODSParameters interface {
	Values() *url.Values
}

//BasicParameters : Basic ODS parameters
type BasicParameters struct {
	Query  string
	Rows   uint
	Format string
}

//Values : Convert parameters to url.Values type
func (p *BasicParameters) Values() *url.Values {
	values := &url.Values{}

	values.Set("q", p.Query)
	values.Set("rows", string(p.Rows))
	values.Set("format", p.Format)

	return values
}

// RecordDownloadParameters : Records Download API parameters
type RecordDownloadParameters struct {
	BasicParameters
}
