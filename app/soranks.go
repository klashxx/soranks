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
	maxErrors = 3
	maxPages  = 1100
	minRep    = 500
)

var (
	location     = flag.String("location", ".", "finder regex")
	jsonfile     = flag.String("json", "", "json sample file (offline)")
	limit        = flag.Int("limit", 20, "max number of records")
	term         = flag.Bool("term", true, "print output to terminal")
	publish      = flag.String("publish", "", "values: 'local' or remote filename (NO ext)")
	author       = lib.Committer{Name: "klasxx", Email: "klashxx@gmail.com"}
	branch       = "dev"
	offline      = true
	stop         = false
	streamErrors = 0
	currentPage  = 1
	lastPage     = currentPage
	counter      = 0
	users        = new(lib.SOUsers)
	key          = ""
	token        = ""
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

	var ranks lib.Ranks
	var err error

	if *jsonfile == "" {
		offline = false
		key, err = lib.GetAPIKey()
		if err == nil {
			key = fmt.Sprintf("&key=%s", key)
		}
	}

	if *publish != "" && *publish != "local" {
		token = lib.GetToken()
	}

	for {
		if offline {
			lib.Info.Println("Extracting from source JSON file.")
			err = lib.StreamFile(*jsonfile, users)
			if err != nil {
				lib.Error.Println("Can't decode json file.")
				os.Exit(5)
			}
			stop = true
		} else {
			lib.Trace.Printf("Requesting page: %d\n", currentPage)

			url := fmt.Sprintf("%s/%s%s", lib.SoAPIURL, fmt.Sprintf(lib.SoUsersQuery, currentPage), key)
			err = lib.StreamHTTP(url, users, true)

			lib.Trace.Printf("Page users: %d\n", len(users.Items))
			if err != nil || len(users.Items) == 0 {

				lib.Warning.Println("Can't stream data.")
				streamErrors++
				if streamErrors >= maxErrors {
					lib.Error.Println("Max retry number reached")
					os.Exit(5)
				}
				continue
			}
		}

		if len(users.Items) == 0 {
			lib.Error.Println("Can't get user info.")
			os.Exit(5)
		}

		lib.Trace.Println("User info extraction.")

		repLimit := lib.GetUserInfo(users, minRep, re, &counter, *limit, &ranks, *term, offline, key)
		if !repLimit {
			break
		}

		lastPage = currentPage
		currentPage++
		if (currentPage >= maxPages && maxPages != 0) || !users.HasMore || stop {
			break
		}
	}

	if counter == 0 {
		lib.Warning.Println("No results found.")
		os.Exit(0)
	}

	if *publish != "" {
		if err := lib.DumpLauncher(ranks, location); err != nil {
			fmt.Println(err)
			os.Exit(5)
		}

		if *publish != "local" {
			if err := lib.GHPublisher(token, publish, branch, author); err != nil {
				fmt.Println(err)
				os.Exit(5)
			}
		}
	}
	lib.Info.Printf("%04d pages requested.\n", lastPage)
	lib.Info.Printf("%04d users found.\n", counter)
}
