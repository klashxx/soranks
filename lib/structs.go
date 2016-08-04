package lib

import (
	"log"
	"time"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
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

type SOUserRank struct {
	Rank         int    `json:"rank"`
	AccountID    int    `json:"account_id"`
	DisplayName  string `json:"display_name"`
	Reputation   int    `json:"reputation"`
	Location     string `json:"location,omitempty"`
	WebsiteURL   string `json:"website_url,omitempty"`
	Link         string `json:"link"`
	ProfileImage string `json:"profile_image"`
}

type Ranks []SOUserRank

type GitHubUpdatePut struct {
	Message   string `json:"message"`
	Committer struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"committer"`
	Content string `json:"content"`
	Sha     string `json:"sha"`
}

type GitHubUpdateRsp struct {
	Content struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		Sha         string `json:"sha"`
		Size        int    `json:"size"`
		URL         string `json:"url"`
		HTMLURL     string `json:"html_url"`
		GitURL      string `json:"git_url"`
		DownloadURL string `json:"download_url"`
		Type        string `json:"type"`
		Links       struct {
			Self string `json:"self"`
			Git  string `json:"git"`
			HTML string `json:"html"`
		} `json:"_links"`
	} `json:"content"`
	Commit struct {
		Sha     string `json:"sha"`
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
		Author  struct {
			Date  time.Time `json:"date"`
			Name  string    `json:"name"`
			Email string    `json:"email"`
		} `json:"author"`
		Committer struct {
			Date  time.Time `json:"date"`
			Name  string    `json:"name"`
			Email string    `json:"email"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			URL string `json:"url"`
			Sha string `json:"sha"`
		} `json:"tree"`
		Parents []struct {
			URL     string `json:"url"`
			HTMLURL string `json:"html_url"`
			Sha     string `json:"sha"`
		} `json:"parents"`
	} `json:"commit"`
}

type GHReqError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

type Repo struct {
	Sha  string `json:"sha"`
	URL  string `json:"url"`
	Tree []struct {
		Path string `json:"path"`
		Mode string `json:"mode"`
		Type string `json:"type"`
		Sha  string `json:"sha"`
		Size int    `json:"size,omitempty"`
		URL  string `json:"url"`
	} `json:"tree"`
	Truncated bool `json:"truncated"`
}

type Create struct {
	Path      string `json:"path"`
	Message   string `json:"message"`
	Content   string `json:"content"`
	Branch    string `json:"branch"`
	Committer `json:"committer"`
}

type Update struct {
	Path      string `json:"path"`
	Message   string `json:"message"`
	Content   string `json:"content"`
	Sha       string `json:"sha"`
	Branch    string `json:"branch"`
	Committer `json:"committer"`
}

type Committer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}