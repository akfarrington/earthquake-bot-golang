package localstorage

import (
	"fmt"
	"log"
	"time"

	"github.com/akfarrington/earthquake_bot/cwbapi"
	"github.com/boltdb/bolt"
)

// StoreLastEqTime saves the lsat earthquake's time in the db
func StoreLastEqTime(db *bolt.DB) {
	bigResponse, smallResponse := cwbapi.GetJSON("")

	var bigTime time.Time
	var smallTime time.Time

	layout := "2006-01-02 15:04:05"
	outputlayout := "2006-01-02T15:04:05"

	// always returns with one quake, but check anyway
	if len(bigResponse.Records.Earthquake) >= 1 {
		bigTime, _ = time.Parse(layout, bigResponse.Records.Earthquake[0].EarthquakeInfo.OriginTime)
		fmt.Println("last big eq is " + bigTime.String())
	}
	if len(smallResponse.Records.Earthquake) >= 1 {
		smallTime, _ = time.Parse(layout, smallResponse.Records.Earthquake[0].EarthquakeInfo.OriginTime)
		fmt.Println("last small eq is " + smallTime.String())
	}

	// add the quake that happened most recently, and add a second so I don't get duplicate quakes
	if smallTime.After(bigTime) {
		StoreTime(db, smallTime.Add(time.Second*1).Format(outputlayout))
	}
	if bigTime.After(smallTime) {
		StoreTime(db, bigTime.Add(time.Second*1).Format(outputlayout))
	}

	new, err := GetTime(db)
	if err != nil {
		log.Println("failed to get new time")
		log.Fatal(err)
	}

	fmt.Println("Saved this time as new \"lasttime\" " + new)
}
