package odsrestclient

import (
	"net/url"
	"strconv"
)

//RecordsParameters : Basic records parameters
type RecordsParameters struct {
	Format      string
	PrettyPrint bool
}

//Values : Convert parameters to url.Values type
func (p RecordsParameters) Values() *url.Values {
	values := &url.Values{}

	values.Set("format", p.Format)
	values.Set("pretty_print", strconv.FormatBool(p.PrettyPrint))

	return values
}

// RecordsDownloadParameters : Records Download API parameters
type RecordsDownloadParameters struct {
	RecordsParameters
	Datasetid string
}

// Values export parameters as url.Values
func (p RecordsDownloadParameters) Values() *url.Values {
	values := p.RecordsParameters.Values()

	values.Set("dataset", p.Datasetid)

	return values
}
