package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
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

	"github.com/klashxx/soranks/lib"
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

var (
	Trace    *log.Logger
	Info     *log.Logger
	Warning  *log.Logger
	Error    *log.Logger
	author   = lib.Committer{Name: "klasxx", Email: "klashxx@gmail.com"}
	branch   = "dev"
	location = flag.String("location", ".", "location")
	jsonfile = flag.String("json", "", "json sample file")
	jsonrsp  = flag.String("jsonrsp", "", "json response file")
	mdrsp    = flag.String("mdrsp", "", "markdown response file")
	limit    = flag.Int("limit", 20, "max number of records")
	term     = flag.Bool("term", false, "print output in terminal")
	publish  = flag.String("publish", "", "publish ranks in Github")
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

func Decode(r io.Reader) (users *lib.SOUsers, err error) {

	users = new(lib.SOUsers)
	return users, json.NewDecoder(r).Decode(users)
}

func StreamHTTP(page int, key string) (users *lib.SOUsers, err error) {

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

func StreamFile(jsonfile string) (users *lib.SOUsers, err error) {
	reader, err := os.Open(jsonfile)
	defer reader.Close()
	return Decode(reader)
}

func GetUserInfo(users *lib.SOUsers, location *regexp.Regexp, counter *int, limit int, ranks *lib.Ranks, term bool) (rep bool) {

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

			s := lib.SOUserRank{Rank: *counter,
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

func DumpJson(path *string, ranks *lib.Ranks) {
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

func DumpMarkdown(path *string, ranks lib.Ranks) {
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

func StreamHTTP2(url string) (repo *lib.Repo, err error) {

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

func Decode2(r io.Reader) (repo *lib.Repo, err error) {

	repo = new(lib.Repo)
	return repo, json.NewDecoder(r).Decode(repo)
}

func Decode3(r io.Reader) (up *lib.GHReqError, err error) {

	up = new(lib.GHReqError)
	return up, json.NewDecoder(r).Decode(up)
}

func Markdown2Base64(path string) (b64 string, err error) {

	mdraw, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Can't load markdown: %s", err)
	}
	return base64.StdEncoding.EncodeToString(mdraw), nil
}

func DataToGihub(data interface{}) (buf io.ReadWriter, err error) {

	buf = new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(data)
	if err != nil {
		return nil, err
	}
	return buf, nil

}

func GitHubIntegration(md string) (err error) {

	encoded, err := Markdown2Base64(*mdrsp)
	if err != nil {
		Error.Println(err)
		os.Exit(5)
	}

	url := fmt.Sprintf("%s%s", GHApiURL, "/git/trees/dev")
	Trace.Printf("Tree url: %s\n", url)

	folder := false
	repo, _ := StreamHTTP2(url)
	for _, file := range repo.Tree {
		if file.Path == "data" {
			url = file.URL
			folder = true
			break
		}
	}
	if !folder {
		fmt.Errorf("Cant't get data folder url")
	}

	Trace.Printf("md: %s\n", md)

	sha := ""
	repo, _ = StreamHTTP2(url)
	for _, file := range repo.Tree {
		if file.Path == md {
			sha = file.Sha
			break
		}
	}

	url = fmt.Sprintf("%s/contents/data/%s", GHApiURL, md)
	Trace.Println(url)

	token := GetKey(GitHubToken)
	if token == "" {
		Error.Println("Can't get github  token!")
		os.Exit(5)
	}
	Info.Printf("token: %s\n", token)

	var buf io.ReadWriter

	if sha == "" {
		Info.Println("Update not detected.")
		data := lib.Create{
			Path:      *mdrsp,
			Message:   "test",
			Content:   encoded,
			Branch:    branch,
			Committer: author}
		buf, _ = DataToGihub(data)
	} else {
		Info.Printf("Update SHA: %s", sha)
		data := lib.Update{
			Path:      *mdrsp,
			Message:   "test",
			Content:   encoded,
			Sha:       sha,
			Branch:    branch,
			Committer: author}
		buf, _ = DataToGihub(data)
	}

	req, err := http.NewRequest("PUT", url, buf)
	if err != nil {
		Trace.Println(err)
	}
	Trace.Println("Sending header.")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		Trace.Println(err)
	}
	Trace.Println("Response.")

	defer response.Body.Close()

	respstring, _ := Decode3(response.Body)

	Trace.Println(respstring)

	return nil
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

	if *publish != "" && *mdrsp == "" {
		Error.Println("Publish requires mdrsp!!")
		os.Exit(5)
	}

	re := regexp.MustCompile(fmt.Sprintf("(?i)%s", *location))

	stop := false
	streamErrors := 0
	currentPage := 1
	lastPage := currentPage
	counter := 0

	var users *lib.SOUsers
	var ranks lib.Ranks

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
		if *publish != "" {
			_ = GitHubIntegration(*publish)
		}
	}

	if *jsonrsp != "" {
		DumpJson(jsonrsp, &ranks)
	}

	Info.Printf("%04d pages requested.\n", lastPage)
	Info.Printf("%04d users found.\n", counter)
}
