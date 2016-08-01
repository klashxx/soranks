package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	MaxErrors     = 3
	MaxPages      = 1
	MinReputation = 400
	APIKeyPath    = "./_secret/api.key"
	GitHubToken   = "./_secret/token"
	SOApiURL      = "https://api.stackexchange.com/2.2/users?page="
	SOQuery       = "pagesize=100&order=desc&sort=reputation&site=stackoverflow"
	GHApiURL      = "https://api.github.com/repos/klashxx/soranks"
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

var (
	Trace    *log.Logger
	Info     *log.Logger
	Warning  *log.Logger
	Error    *log.Logger
	location = flag.String("location", ".", "location")
	jsonfile = flag.String("json", "", "json sample file")
	jsonrsp  = flag.String("jsonrsp", "", "json response file")
	mdrsp    = flag.String("mdrsp", "", "markdown response file")
	limit    = flag.Int("limit", 20, "max number of records")
	term     = flag.Bool("term", false, "print output in terminal")
	publish  = flag.Bool("publish", false, "publish ranks in Github")
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

func Decode(r io.Reader) (users *SOUsers, err error) {

	users = new(SOUsers)
	return users, json.NewDecoder(r).Decode(users)
}

func StreamHTTP(page int, key string) (users *SOUsers, err error) {

	var reader io.ReadCloser

	url := fmt.Sprintf("%s%d&%s%s", SOApiURL, page, SOQuery, key)
	Trace.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Trace.Println(err)
		return users, err
	}
	Trace.Println("Sending header.")

	req.Header.Set("Accept-Encoding", "gzip")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		Trace.Println(err)
		return users, err
	}
	Trace.Println("Response.")

	defer response.Body.Close()

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			Trace.Println(err)
			return users, err
		}
		defer reader.Close()
	default:
		reader = response.Body
	}
	return Decode(reader)
}

func StreamFile(jsonfile string) (users *SOUsers, err error) {
	reader, err := os.Open(jsonfile)
	defer reader.Close()
	return Decode(reader)
}

func GetUserInfo(users *SOUsers, location *regexp.Regexp, counter *int, limit int, ranks *Ranks, term bool) (rep bool) {

	for _, user := range users.Items {
		Trace.Printf("Procesing user: %d\n", user.AccountID)
		if user.Reputation < MinReputation {
			return false
		}
		if location.MatchString(user.Location) {
			*counter += 1
			if *counter == 1 && term {
				Info.Println("User data:")
				Info.Printf("%4s %-30s %6s %s\n", "Rank", "Name", "Rep", "Location")
			}

			s := SOUserRank{Rank: *counter,
				AccountID:    user.AccountID,
				DisplayName:  user.DisplayName,
				Reputation:   user.Reputation,
				Location:     user.Location,
				WebsiteURL:   user.WebsiteURL,
				Link:         user.Link,
				ProfileImage: user.ProfileImage}

			*ranks = append(*ranks, s)

			if term {
				Info.Printf("%4d %-30s %6d %s\n", *counter, html.UnescapeString(user.DisplayName),
					user.Reputation, html.UnescapeString(user.Location))
			}

			if *counter >= limit && limit != 0 {
				return false
			}

		}
	}
	return true
}

func DumpJson(path *string, ranks *Ranks) {
	Trace.Printf("Writing JSON to: %s\n", *path)
	jsonenc, _ := json.MarshalIndent(*ranks, "", " ")
	f, err := os.Create(*path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	n4, err := w.WriteString(string(jsonenc))
	if err != nil {
		panic(err)
	}
	Trace.Printf("Wrote %d bytes to %s\n", n4, *path)

	w.Flush()
}

func DumpMarkdown(path *string, ranks Ranks) {
	Trace.Printf("Writing MD to: %s\n", *path)

	head := `# soranks

[Stackoverflow](http://stackoverflow.com/) rankings by **location**.

### Area%s


Rank|Name|Rep|Location|Web|Avatar
----|----|---|--------|---|------
`
	var fmtLocation string

	if *location == "." {
		fmtLocation = ": WorldWide"
	} else {
		fmtLocation = fmt.Sprintf(" *pattern*: %s", *location)
	}

	userfmt := "{{.Rank}}|[{{.DisplayName}}]({{.Link}})|{{.Reputation}}|{{.Location}}|{{.WebsiteURL}}|![Avatar]({{.ProfileImage}})\n"

	f, err := os.Create(*path)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	w := bufio.NewWriter(f)
	n4, err := w.WriteString(fmt.Sprintf(head, fmtLocation))
	if err != nil {
		panic(err)
	}
	w.Flush()

	tmpl, _ := template.New("Ranking").Parse(userfmt)
	for _, userRank := range ranks {
		_ = tmpl.Execute(f, userRank)
	}
	Trace.Printf("Wrote %d bytes to %s\n", n4, *path)
	w.Flush()
}

func GetKey(path string) (key string) {

	_, err := os.Stat(path)
	if err != nil {
		Warning.Printf("Can't find key: %s", path)
		return ""
	}

	strkey, err := ioutil.ReadFile(path)
	if err != nil {
		Warning.Printf("Can't load key: %s", err)
		return ""
	}

	return strings.TrimRight(string(strkey)[:], "\n")
}

func StreamHTTP2(url string) (repo *Repo, err error) {

	Trace.Println(url)

	var reader io.ReadCloser
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Trace.Println(err)
	}
	Trace.Println("Sending header.")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		Trace.Println(err)
	}
	Trace.Println("Response.")

	defer response.Body.Close()

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			Trace.Println(err)
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	return Decode2(reader)
}

func Decode2(r io.Reader) (repo *Repo, err error) {

	repo = new(Repo)
	return repo, json.NewDecoder(r).Decode(repo)
}

func GitHubIntegration() {
	//_ = GetKey(GitHubToken)

	url := fmt.Sprintf("%s%s", GHApiURL, "/git/trees/dev")

	repo, _ := StreamHTTP2(url)
	for _, file := range repo.Tree {
		if file.Path == "data" {
			url = file.URL
			break
		}
	}

	var md string
	if *location == "." {
		md = "global.md"
	} else {
		md = fmt.Sprintf("%s.md", *location)
	}
	Trace.Printf("base md: %s\n", md)

	sha := ""
	repo, _ = StreamHTTP2(url)
	for _, file := range repo.Tree {
		if file.Path == md {
			sha = file.Sha
			break
		}
	}

	if sha == "" {
		Info.Println("Update not detected.")
	} else {
		Info.Printf("Update SHA: %s", sha)
	}
}

func main() {
	flag.Parse()

	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	Trace.Println("location: ", *location)
	Trace.Println("json: ", *jsonfile)
	Trace.Println("jsontest: ", *jsonfile)
	Trace.Println("jsonrsp: ", *jsonrsp)
	Trace.Println("mdrsp: ", *mdrsp)
	Trace.Println("limit: ", *limit)
	Trace.Println("term: ", *term)
	Trace.Println("publish: ", *publish)

	re := regexp.MustCompile(fmt.Sprintf("(?i)%s", *location))

	stop := false
	streamErrors := 0
	currentPage := 1
	lastPage := currentPage
	counter := 0

	var users *SOUsers
	var ranks Ranks

	for {
		if *jsonfile == "" {
			var key string
			var err error
			if lastPage == currentPage {
				Info.Println("Trying to extract API key.")
				key = fmt.Sprintf("&key=%s", GetKey(APIKeyPath))
			}

			Trace.Printf("Requesting page: %d\n", currentPage)

			users, err = StreamHTTP(currentPage, key)

			Trace.Printf("Page users: %d\n", len(users.Items))
			if err != nil || len(users.Items) == 0 {

				Warning.Println("Can't stream data.")
				streamErrors += 1
				if streamErrors >= MaxErrors {
					Error.Println("Max retry number reached")
					os.Exit(5)
				}
				continue
			}
		} else {
			Info.Println("Extracting from source JSON file.")
			var err error
			users, err = StreamFile(*jsonfile)
			if err != nil {
				Error.Println("Can't decode json file.")
				os.Exit(5)
			}
			stop = true
		}

		Trace.Println("User info extraction.")

		repLimit := GetUserInfo(users, re, &counter, *limit, &ranks, *term)
		if !repLimit {
			break
		}
		Trace.Println("User info extraction done.")

		lastPage = currentPage
		currentPage += 1
		if (currentPage >= MaxPages && MaxPages != 0) || !users.HasMore || stop {
			break
		}
	}

	if counter == 0 {
		Warning.Println("No results found.")
		os.Exit(0)
	}

	if *mdrsp != "" {
		DumpMarkdown(mdrsp, ranks)
		if *publish {
			GitHubIntegration()
		}
	}

	if *jsonrsp != "" {
		DumpJson(jsonrsp, &ranks)
	}

	Info.Printf("%04d pages requested.\n", lastPage)
	Info.Printf("%04d users found.\n", counter)
}
