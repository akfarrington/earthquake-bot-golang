package localstorage

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

const dbfilename string = "local.db"
const bucketname string = "vars"
const timekey string = "lasttime"

// OpenDB opens a db and returns the db object
func OpenDB() *bolt.DB {
	db, err := bolt.Open(dbfilename, 0600, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return db
}

// StoreTime stores a time value in the persistent storage
func StoreTime(db *bolt.DB, time string) error {
	// store some data
	bucketName := []byte(bucketname)
	key := []byte(timekey)
	value := []byte(time)

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// GetTime gets the last time variable
func GetTime(db *bolt.DB) (string, error) {
	bucketName := []byte(bucketname)
	key := []byte(timekey)

	var val string

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", bucketName)
		}

		val = string(bucket.Get(key))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return val, nil
}

// CheckDBExists checks if the db file exists (main will add a value if not)
func CheckDBExists() bool {
	info, err := os.Stat(dbfilename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// StoreEqLastTime checks if a new earthquake's time is after the old db's lasttime value,
// then saves it if so (this is in case there are multiple earthquakes)
func StoreEqLastTime(db *bolt.DB, newtime string) {
	inputlayout := "2006-01-02 15:04:05"
	outputlayout := "2006-01-02T15:04:05"

	eqTime, err := time.Parse(inputlayout, newtime)
	if err != nil {
		log.Println("error parsing time from cwb api (StoreEqLastTime)")
		log.Fatal(err)
	}

	lasttimestring, _ := GetTime(db)

	lasttime, err := time.Parse(outputlayout, lasttimestring)
	if err != nil {
		log.Println("error parsing time from db (StoreEqLastTime)")
		log.Fatal(err)
	}

	if eqTime.After(lasttime) || eqTime.Equal(lasttime) {
		log.Println("storing earthquake's time as new lasttime")
		StoreTime(db, eqTime.Add(time.Second*1).Format(outputlayout))
	}
}
