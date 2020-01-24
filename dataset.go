package odsrestclient

import "net/url"

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

// DatasetSearchParameters : Records Download API parameters
type DatasetSearchParameters struct {
	BasicParameters
	Start uint
}

func (p *DatasetSearchParameters) Values() *url.Values {
	return p.BasicParameters.Values()
}
