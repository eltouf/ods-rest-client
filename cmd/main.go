package main

import (
	"context"
	"flag"
	"log"
	"net/url"
	"os"
	"time"

	firebase "firebase.google.com/go"
	odsrestclient "github.com/eltouf/ods-rest-client"
	"google.golang.org/api/option"
)

func main() {

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error initializing user home: %v\n", err)
	}

	opt := option.WithCredentialsFile(home + string(os.PathSeparator) + ".config" + string(os.PathSeparator) + "bmtimetable-381d0e19e193.json")

	config := &firebase.Config{ProjectID: "bmtimetable", StorageBucket: "bmtimetable.appspot.com"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	ctx := context.Background()
	storage, err := app.Storage(ctx)
	if err != nil {
		log.Fatalf("error initializing storage: %v\n", err)
	}

	bkt, err := storage.DefaultBucket()
	if err != nil {
		log.Fatalf("error initializing bucket: %v\n", err)
	}

	query := flag.String("query", "saeiv", "Query term")

	parameters := odsrestclient.DatasetSearchParameters{
		DatasetParameters: odsrestclient.DatasetParameters{
			Format:      "json",
			PrettyPrint: false,
		},
		Start: 0,
		Query: *query,
		Rows:  50,
	}

	baseURL, err := url.Parse("https://opendata.bordeaux-metropole.fr/")
	if err != nil {
		log.Panic(err)
	}

	odsClient := odsrestclient.NewClient(nil, baseURL)
	catalog, err := odsClient.DatasetSearch(parameters)

	for _, dataset := range catalog.Datasets {

		file := bkt.Object("data/raw/" + time.Now().Format("2006/01/02") + "/" + dataset.Datasetid).NewWriter(ctx)

		log.Printf("object name %v\n", file.Name)

		odsClient.DownloadRecords(
			odsrestclient.RecordsDownloadParameters{
				RecordsParameters: odsrestclient.RecordsParameters{
					Format:      "json",
					PrettyPrint: false,
				},
				Datasetid: dataset.Datasetid,
			},
			file,
		)

		if err := file.Close(); err != nil {
			log.Fatalf("error closing file : %v\n", err)
		}

		log.Printf("object name %v\n", file.Size)
	}
}
