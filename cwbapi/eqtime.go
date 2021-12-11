package cwbapi

import "time"

// GetTwTime returns a tw time string (that the api can use)
func GetTwTime() string {
	loc, _ := time.LoadLocation("Asia/Taipei")
	now := time.Now().In(loc).Format("2006-01-02T15:04:05")
	return now
}
