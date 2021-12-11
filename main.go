package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akfarrington/earthquake_bot/cwbapi"
	"github.com/akfarrington/earthquake_bot/eqlog"
	"github.com/akfarrington/earthquake_bot/localstorage"
	"github.com/akfarrington/earthquake_bot/twitter"

	"github.com/ChimeraCoder/anaconda"
	"github.com/boltdb/bolt"
	"github.com/joho/godotenv"
)

func main() {
	// change logger to use taiwan time
	log.SetFlags(0)
	log.SetOutput(new(eqlog.LogWriter))

	// get environment vars
	err := godotenv.Load()
	if err != nil {
		log.Print("loading .env file failed")
		log.Fatal(err)
	}

	// delete this
	if os.Getenv("ENVIRONMENT") == "PROD" {
		log.Println("******** production ********")
	} else {
		log.Println("******** development ********")
	}

	//first check if localstorage file exists
	if localstorage.CheckDBExists() == false {
		// no file, create one
		log.Printf("no db file exists - creating one now")
		db := localstorage.OpenDB()

		// store last eq data
		localstorage.StoreLastEqTime(db)
		log.Print("local storage file last eq time, exiting now")
		db.Close()
		os.Exit(1)
	}

	// open connection to the database
	db := localstorage.OpenDB()
	defer db.Close()

	// get twitter api ready
	tapi := twitter.ConnectAPI()

	// get last time from db
	lasttime, err := localstorage.GetTime(db)
	if err != nil {
		log.Fatal("couldn't open db")
	}

	log.Println("Started - using old lasttime: " + lasttime)

	// handle ctrl z
	SetupCloseHandler(db, tapi)

	// print this to make sure I know it passed all the previous shit
	fmt.Println("Starting loop now")

	for i := 1; ; i++ {
		lasttime, err = localstorage.GetTime(db)
		if err != nil {
			log.Fatal("can't get lasttime from database")
		}
		// get json data from the api server
		respBig, respSmall := cwbapi.GetJSON(lasttime)

		// process the json data
		processEq(respBig, tapi, db)
		processEq(respSmall, tapi, db)

		// use counter to update every ~ hour
		if i == 1 {
			log.Print("completed once")
		} else if i%60 == 0 {
			log.Print("still running ~1 hour, no hangs")
		}

		// wait for 1 minute
		time.Sleep(1 * time.Minute)
	}
}

func processEq(eq cwbapi.Response, tapi *anaconda.TwitterApi, db *bolt.DB) {
	for _, earthquake := range eq.Records.Earthquake {
		tweetString := earthquake.ReportContent + " #台灣 #地震 #taiwan #earthquake"
		log.Println(tweetString)
		if os.Getenv("ENVIRONMENT") == "PROD" {
			twitter.Post(tapi, tweetString, earthquake.ReportImageURI, eq)
		} else {
			fmt.Println(tweetString)
		}
		//store earthquake as last time
		localstorage.StoreEqLastTime(db, earthquake.EarthquakeInfo.OriginTime)
	}
}

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
// copied from https://golangcode.com/handle-ctrl-c-exit-in-terminal/
func SetupCloseHandler(db *bolt.DB, tapi *anaconda.TwitterApi) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		db.Close()
		tapi.Close()
		os.Exit(0)
	}()
}
