package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/klashxx/soranks/lib"
)

const (
	MaxErrors     = 3
	MaxPages      = 1100
	MinReputation = 500
)

var (
	author   = lib.Committer{Name: "klasxx", Email: "klashxx@gmail.com"}
	branch   = "dev"
	location = flag.String("location", ".", "finder regex")
	jsonfile = flag.String("json", "", "json sample file (offline)")
	limit    = flag.Int("limit", 20, "max number of records")
	term     = flag.Bool("term", false, "print output to terminal")
	publish  = flag.String("publish", "", "values: 'local' or remote md filename.")
)

func main() {
	flag.Parse()

	lib.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	lib.Trace.Println("location: ", *location)
	lib.Trace.Println("json: ", *jsonfile)
	lib.Trace.Println("jsontest: ", *jsonfile)
	lib.Trace.Println("limit: ", *limit)
	lib.Trace.Println("term: ", *term)
	lib.Trace.Println("publish: ", *publish)

	re := regexp.MustCompile(fmt.Sprintf("(?i)%s", *location))

	stop := false
	streamErrors := 0
	currentPage := 1
	lastPage := currentPage
	counter := 0

	var ranks lib.Ranks
	var key string
	var err error

	users := new(lib.SOUsers)

	for {
		if *jsonfile == "" {
			if lastPage == currentPage {
				lib.Info.Println("Trying to extract API key.")
				key = fmt.Sprintf("&key=%s", lib.GetKey(lib.APIKeyPath))
			}

			lib.Trace.Printf("Requesting page: %d\n", currentPage)

			url := fmt.Sprintf("%s/%s%s", lib.SOApiURL, fmt.Sprintf(lib.SOUsersQuery, currentPage), key)
			err = lib.StreamHTTP(url, users, true)

			lib.Trace.Printf("Page users: %d\n", len(users.Items))
			if err != nil || len(users.Items) == 0 {

				lib.Warning.Println("Can't stream data.")
				streamErrors += 1
				if streamErrors >= MaxErrors {
					lib.Error.Println("Max retry number reached")
					os.Exit(5)
				}
				continue
			}
		} else {
			lib.Info.Println("Extracting from source JSON file.")
			err = lib.StreamFile(*jsonfile, users)
			if err != nil {
				lib.Error.Println("Can't decode json file.")
				os.Exit(5)
			}
			stop = true
		}

		if len(users.Items) == 0 {
			lib.Error.Println("Can't get user info.")
			os.Exit(5)
		}

		lib.Trace.Println("User info extraction.")

		repLimit := lib.GetUserInfo(users, MinReputation, re, &counter, *limit, &ranks, *term)
		if !repLimit {
			break
		}

		lastPage = currentPage
		currentPage += 1
		if (currentPage >= MaxPages && MaxPages != 0) || !users.HasMore || stop {
			break
		}
	}

	if counter == 0 {
		lib.Warning.Println("No results found.")
		os.Exit(0)
	}

	if *publish != "" {

		if err = lib.DumpJson(&ranks); err != nil {
			lib.Error.Println("JSON Dump failed:", err)
			os.Exit(5)
		}
		if err = lib.DumpMarkdown(ranks, location); err != nil {
			lib.Error.Println("MD Dump failed:", err)
			os.Exit(5)
		}

		if *publish != "local" {
			if err = lib.GitHubConnector(*publish, branch, author); err != nil {
				lib.Error.Println("GitHub connection error:", err)
				os.Exit(5)
			}
		}
	}

	lib.Info.Printf("%04d pages requested.\n", lastPage)
	lib.Info.Printf("%04d users found.\n", counter)
}
