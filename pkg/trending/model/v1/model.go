package model

import (
	"context"
	"log"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

const orgName = "test"
const bucketName = "github-trending"

type Project struct {
	Name      string `json:"name,omitempty"`
	Url       string `json:"url,omitempty"`
	Overview  string `json:"overview,omitempty"`
	Star      int    `json:"star,omitempty"`
	TodayStar int    `json:"todayStar,omitempty"`
	Fork      int    `json:"fork,omitempty"`
}

func Create(client influxdb.Client) error {
	organizationAPI := client.OrganizationsAPI()
	organization, err := organizationAPI.FindOrganizationByName(context.TODO(), orgName)
	if err != nil {
		return err
	}

	bucketsAPI := client.BucketsAPI()

	if bucket, err := bucketsAPI.FindBucketByName(context.TODO(), bucketName); bucket != nil && err == nil {
		log.Printf("bucket %s has exists", bucketName)
		return nil
	}

	if _, err := bucketsAPI.CreateBucketWithName(context.TODO(), organization, bucketName); err != nil {
		return err
	}

	return nil
}

func SaveDailyTrending(client influxdb.Client, date time.Time, data []*Project) error {
	writeAPI := client.WriteAPI(orgName, bucketName)

	for _, v := range data {
		p := influxdb.NewPointWithMeasurement("daily").
			AddTag("name", v.Name).
			AddTag("url", v.Url).
			AddTag("overview", v.Overview).
			AddField("star", v.Star).
			AddField("todayStar", v.TodayStar).
			AddField("fork", v.Fork).
			SetTime(date)

		writeAPI.WritePoint(p)
	}

	writeAPI.Flush()
	return nil
}
