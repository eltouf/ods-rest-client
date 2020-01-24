package main

import odsrestclient "github.com/eltouf/ods-rest-client"

func main() {
	parameters := &odsrestclient.DatasetSearchParameters{
		BasicParameters: odsrestclient.BasicParameters{
			Query:  "saiev",
			Rows:   10,
			Format: "json",
		},
		Start: 0,
	}

	odsrestclient.NewClient(nil).DatasetSearch(parameters)
}
