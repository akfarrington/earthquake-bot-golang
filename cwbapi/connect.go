package cwbapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const eqRoot string = "https://opendata.cwb.gov.tw/api/v1/rest/datastore/"
const eqBig string = "E-A0015-001"
const eqSmall string = "E-A0016-001"

// GetJSON got this info from https://blog.alexellis.io/golang-json-api-client/
func GetJSON(lasttime string) (Response, Response) {
	// get api authentication code
	authCode := os.Getenv("CWB_AUTH_CODE")

	// no "last time" if lasttime is empty (for new local DBs)
	bigURL := eqRoot + eqBig + "?Authorization=" + authCode
	smallURL := eqRoot + eqSmall + "?Authorization=" + authCode

	// while bot is running normally, this will add the last time to make sure
	// I don't get duplicate earthquakes
	if len(lasttime) > 0 {
		bigURL = bigURL + "&timeFrom=" + lasttime
		smallURL = smallURL + "&timeFrom=" + lasttime
	}

	bigResp := fetchJSON(bigURL)
	time.Sleep(time.Second * 1)
	smallResp := fetchJSON(smallURL)

	return bigResp, smallResp
}

func fetchJSON(url string) Response {
	// struct to put useful json data in
	response := Response{}

	if os.Getenv("ENVIRONMENT") == "DEV" {
		log.Println(url)
	}

	// this hopefully fixes some bullshit problems with the api
	client := http.Client{
		Timeout: time.Second * 10,
		Transport: (&http.Transport{
			Dial:                (&net.Dialer{Timeout: 10 * time.Second}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		}),
	}

	// this one has errored out many times
	// so it'll loop to attempt again after two minutes, attempting 10 times
	res, getErr := client.Get(url)
	for i := 0; i < 10; i++ {
		if getErr != nil {
			log.Print("failed to fetch json data from api, trying again")
			time.Sleep(2 * time.Minute)

			// after waiting 2 minutes, trying to get the url again
			res, getErr = client.Get(url)
		} else {
			break
		}
	}

	// returning an empty response if the above fails too many times
	// hopefully this won't happen
	if getErr != nil {
		log.Println("failed too many times, so breaking for now")
		return response
	}

	// in case it continues to fail, return an empty result
	// (this resulted in a null pointer exception last time, so hopefully this
	// fixes it)
	if getErr != nil {
		return response
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println("error reading CWB response body")
		log.Println(readErr)
	}

	jsonErr := json.Unmarshal(body, &response)
	for i := 0; i < 3; i++ {
		if jsonErr != nil {
			log.Print("error unmarshaling the CWB json (probably maintenance), will wait for an hour")
			log.Println(jsonErr)
			time.Sleep(1 * time.Hour)
			// give up and return empty response struct
			return Response{}
		}
		break
	}

	return response
}
