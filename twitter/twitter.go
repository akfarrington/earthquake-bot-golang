package twitter

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/akfarrington/earthquake_bot/cwbapi"
	"github.com/akfarrington/earthquake_bot/taiwanmap"
)

// ConnectAPI connects to the api and save the api object
func ConnectAPI() *anaconda.TwitterApi {
	apiKey := os.Getenv("API_KEY")
	apiSecretKey := os.Getenv("API_SECRET_KEY")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")

	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, apiKey, apiSecretKey)
	return api
}

// Post posts everything sent to it including the tweet text and the image found at the url
func Post(tapi *anaconda.TwitterApi, text string, imageURL string, response cwbapi.Response) {
	// try to get the image from cwb's server
	// cwbimage := base64.StdEncoding.EncodeToString([]byte(getImage(imageURL)))

	// the image isn't loading so make my own

	coordlist := []taiwanmap.Coords{}

	for _, area := range response.Records.Earthquake[0].Intensity.ShakingArea {
		for _, station := range area.EqStation {
			newcoord := taiwanmap.NewCoords(station.StationLat.Value, station.StationLong.Value, station.StationIntensity.Value)
			coordlist = append(coordlist, newcoord)
		}
	}

	pic := taiwanmap.GetBasePic()

	// mark the epicenter, then color code stations that felt the quake according to the key
	tmap := taiwanmap.MarkEpicenter(pic, taiwanmap.NewCoords(response.Records.Earthquake[0].EarthquakeInfo.Epicenter.Long.Value, response.Records.Earthquake[0].EarthquakeInfo.Epicenter.Lat.Value, 1))
	tmap = taiwanmap.MarkPicList(tmap, coordlist)

	buf := new(bytes.Buffer)

	png.Encode(buf, tmap)
	imageBit := buf.Bytes()

	cwbimage := base64.StdEncoding.EncodeToString([]byte(imageBit))

	// upload the media
	image, err := tapi.UploadMedia(cwbimage)

	// put the media id string into this string map
	urlValues := make(map[string][]string)
	urlValues["media_ids"] = []string{image.MediaIDString}

	// post the tweet text with media ids
	tweet, err := tapi.PostTweet(text, urlValues)
	if err != nil {
		log.Print("failed to post the tweet")
		log.Fatal(err)
	}

	fmt.Println("Tweet successful - ID: " + tweet.IdStr)
}

func getImage(url string) string {
	// get default stuff for custom transport
	defaultTransport := http.DefaultTransport.(*http.Transport)

	// Create new Transport that ignores self-signed SSL
	// (previously, CWB used a self-signed cert for the subdomain where the images were stored)
	customTransport := &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		// set this to a much higher number than necessary
		// bot is running in another country, and for some reason
		// takes well over 10 seconds for wget to download the image
		// so I'm making this 1 minute to be safe (should be more than enough time)
		// plus who cares if it takes longer to download the next image
		Timeout:   time.Minute * 1,
		Transport: customTransport, // finishes removing the security check
	}

	// do the request
	// if there's a big earthquake, the cwb servers are overloaded and can't serve images
	// and I get a tls handshake error. Keep trying until it works, or just fail forever
	// if it takes 20 minutes
	res, getErr := client.Get(url)
	if getErr != nil {
		log.Print("failed to do/complete request for eq image, trying again in 2 minutes")
		return ""
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Print("failed to read response for eq image")
		log.Fatal(readErr)
	}

	return string(body)
}
