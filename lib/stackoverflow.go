package lib

import (
	"fmt"
	"html"
	"os"
	"regexp"
)

func GetAPIKey() (key string, err error) {
	Info.Println("Trying to load API key.")
	key = os.Getenv("API_KEY")
	if key != "" {
		return key, nil
	}
	Warning.Println("Can't load API_KEY env variable")

	Info.Println("Trying to extract API key.")
	key, err = GetKey(APIKeyPath)
	if err != nil {
		Warning.Println(err)
		return "", err
	}
	return key, nil
}

func GetTags(userid int, key string, offline bool) string {

	if offline {
		Trace.Printf("Not avaliable offline")
		return ""
	}

	tags := new(SOTopTags)
	tagstr := ""

	url := fmt.Sprintf("%s/%s%s", SOApiURL, fmt.Sprintf(SOUserTags, userid), key)

	if err := StreamHTTP(url, tags, true); err != nil {
		Trace.Printf("No info for: %d\n", userid)
		return ""
	}

	if len(tags.Items) == 0 {
		Trace.Printf("No info for: %d\n", userid)
		return ""
	}

	for _, tag := range tags.Items {
		tagstr = fmt.Sprintf("%s<li>%s</li>", tagstr, tag.TagName)
	}

	return tagstr
}

func GetUserInfo(users *SOUsers, min int, location *regexp.Regexp, counter *int, limit int, ranks *Ranks, term bool, offline bool, key string) bool {

	for _, user := range users.Items {
		Trace.Printf("Procesing user: %d\n", user.UserID)

		if user.Reputation < min {
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
				ProfileImage: user.ProfileImage,
				TopTags:      GetTags(user.UserID, key, offline)}

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
