package odsrestclient

import "net/url"

import "strconv"

//Dataset dataset
type Dataset struct {
	Datasetid string
	Metas     struct {
		Publisher         string
		Domain            string
		RecordsCount      uint `json:"records_count"`
		Title             string
		MetadataProcessed string `json:"metadata_processed"`
		DataProcessed     string `json:"data_processed"`
	}
	HasRecords bool
	fields     []interface{}
}

//Catalog : A set a dataset
type Catalog struct {
	Result
	Datasets []Dataset
}

//DatasetParameters : Dataset Parameters
type DatasetParameters struct {
	Format      string
	PrettyPrint bool
}

//Values : Convert parameters to url.Values type
func (p DatasetParameters) Values() *url.Values {
	values := &url.Values{}

	values.Set("format", p.Format)
	values.Set("pretty_print", strconv.FormatBool(p.PrettyPrint))

	return values
}

// DatasetSearchParameters : Records Download API parameters
type DatasetSearchParameters struct {
	DatasetParameters
	Query string
	Rows  uint
	Start uint
}

// Values export parameters as url.Values
func (p DatasetSearchParameters) Values() *url.Values {
	values := p.DatasetParameters.Values()

	values.Set("start", strconv.FormatUint(uint64(p.Start), 10))
	values.Set("q", p.Query)
	values.Set("rows", strconv.FormatUint(uint64(p.Rows), 10))

	return values
}
