package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	ApiURL = "https://api.stackexchange.com/2.2/users"
	Query  = "pagesize=100&order=desc&sort=reputation&site=stackoverflow"
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

func Decode(r io.Reader) (x *SOUsers, err error) {
	x = new(SOUsers)
	err = json.NewDecoder(r).Decode(x)
	return x, err
}

func main() {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?page=%d&%s", ApiURL, 1, Query), nil)
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

	users, err := Decode(reader)

	for _, user := range users.Items {
		fmt.Println(user.DisplayName)
		fmt.Println(user.Reputation)
		fmt.Println(user.Location)
	}

}
