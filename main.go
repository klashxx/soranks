package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	MaxPages      = 1
	MinReputation = 200
	APIKeyPath    = "./_secret/api.key"
	ApiURL        = "https://api.stackexchange.com/2.2/users"
	CQuery        = "pagesize=100&order=desc&sort=reputation&site=stackoverflow"
)

type SOUsers struct {
	Items []struct {
		BadgeCounts struct {
			Bronze int `json:"bronze"`
			Silver int `json:"silver"`
			Gold   int `json:"gold"`
		} `json:"badge_counts"`
		AccountID               int    `json:"account_id"`
		IsEmployee              bool   `json:"is_employee"`
		LastModifiedDate        int    `json:"last_modified_date"`
		LastAccessDate          int    `json:"last_access_date"`
		Age                     int    `json:"age,omitempty"`
		ReputationChangeYear    int    `json:"reputation_change_year"`
		ReputationChangeQuarter int    `json:"reputation_change_quarter"`
		ReputationChangeMonth   int    `json:"reputation_change_month"`
		ReputationChangeWeek    int    `json:"reputation_change_week"`
		ReputationChangeDay     int    `json:"reputation_change_day"`
		Reputation              int    `json:"reputation"`
		CreationDate            int    `json:"creation_date"`
		UserType                string `json:"user_type"`
		UserID                  int    `json:"user_id"`
		AcceptRate              int    `json:"accept_rate,omitempty"`
		Location                string `json:"location,omitempty"`
		WebsiteURL              string `json:"website_url,omitempty"`
		Link                    string `json:"link"`
		ProfileImage            string `json:"profile_image"`
		DisplayName             string `json:"display_name"`
	} `json:"items"`
	HasMore        bool `json:"has_more"`
	QuotaMax       int  `json:"quota_max"`
	QuotaRemaining int  `json:"quota_remaining"`
}

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARN: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func Decode(r io.Reader) (x *SOUsers, err error) {
	x = new(SOUsers)
	return x, json.NewDecoder(r).Decode(x)
}

func Streamdata(page int, key string) (x *SOUsers, err error) {

	url := fmt.Sprintf("%s?page=%d&%s%s", ApiURL, page, CQuery, key)

	Trace.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error 1", err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error 2", err)
	}
	defer response.Body.Close()

	var reader io.ReadCloser

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}
	return Decode(reader)
}

func GetKey() (key string, err error) {
	_, err = os.Stat(APIKeyPath)

	if err != nil {
		return "", fmt.Errorf("Can't find API key: %s", APIKeyPath)
	}

	strkey, err := ioutil.ReadFile(APIKeyPath)
	if err != nil {
		return "", fmt.Errorf("Can't load API key: %s", err)
	}

	return fmt.Sprintf("&key=%s", strings.TrimRight(string(strkey)[:], "\n")), nil
}

func main() {

	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	stop := false
	currentPage := 1
	key, err := GetKey()

	if err != nil {
		Warning.Println(err)
	}

	for {
		users, _ := Streamdata(currentPage, key)

		for _, user := range users.Items {
			fmt.Println(user.DisplayName)
			fmt.Println(user.Reputation)
			fmt.Println(user.Location)
			if user.Reputation < MinReputation {
				stop = true
			}
		}

		currentPage += 1
		if currentPage >= MaxPages || !users.HasMore || stop {
			break
		}
	}

}
