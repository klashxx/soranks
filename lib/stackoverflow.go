package lib

import (
	"fmt"
	"html"
	"regexp"
)

func GetUserInfo(users *SOUsers, min int, location *regexp.Regexp, counter *int, limit int, ranks *Ranks, term bool, offline bool, key string) bool {

	var url string
	var tagstr string
	tags := new(SOTopTags)

	for _, user := range users.Items {
		tagstr = ""
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

			if !offline {
				url = fmt.Sprintf("%s/%s%s", SOApiURL, fmt.Sprintf(SOUserTags, user.UserID), key)

				if err := StreamHTTP(url, tags, true); err != nil {
					Trace.Printf("No info for: %d\n", user.UserID)
				} else {
					if len(tags.Items) > 0 {
						for _, tag := range tags.Items {
							tagstr = fmt.Sprintf("%s<li>%s</li>", tagstr, tag.TagName)
						}
					} else {
						Trace.Printf("No info for: %d\n", user.UserID)
					}
				}

			}

			s := SOUserRank{Rank: *counter,
				AccountID:    user.AccountID,
				DisplayName:  user.DisplayName,
				Reputation:   user.Reputation,
				Location:     user.Location,
				WebsiteURL:   user.WebsiteURL,
				Link:         user.Link,
				ProfileImage: user.ProfileImage,
				TopTags:      tagstr}

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
