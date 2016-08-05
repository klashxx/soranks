package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func GitHubConnector(apiurl string, md string, mdfmt string, tokenpath string, branch string, author Committer) (err error) {

	encoded, err := Markdown2Base64(mdfmt)
	if err != nil {
		Error.Println(err)
		os.Exit(5)
	}

	url := fmt.Sprintf("%s%s", apiurl, "/git/trees/dev")
	Trace.Printf("Tree url: %s\n", url)

	folder := false
	repo := new(Repo)
	_ = StreamHTTP(url, repo, false)
	for _, file := range repo.Tree {
		if file.Path == "data" {
			url = file.URL
			folder = true
			break
		}
	}
	if !folder {
		return fmt.Errorf("Cant't get data folder url")
	}

	Trace.Printf("md: %s\n", md)

	sha := ""
	_ = StreamHTTP(url, repo, false)
	for _, file := range repo.Tree {
		if file.Path == md {
			sha = file.Sha
			break
		}
	}

	url = fmt.Sprintf("%s/contents/data/%s", apiurl, md)
	Trace.Println(url)

	token := GetKey(tokenpath)
	if token == "" {
		Error.Println("Can't get github  token!")
		os.Exit(5)
	}
	Info.Printf("token: %s\n", token)

	var buf io.ReadWriter
	c := fmt.Sprintf("%s [%s]", md, time.Now().Format(time.RFC3339))

	if sha == "" {
		Info.Println("Update not detected.")
		data := Create{
			Path:      mdfmt,
			Message:   fmt.Sprintf("Create: %s", c),
			Content:   encoded,
			Branch:    branch,
			Committer: author}
		buf, _ = DataToGihub(data)
	} else {
		Info.Printf("Update SHA: %s", sha)
		data := Update{
			Path:      mdfmt,
			Message:   fmt.Sprintf("Update: %s", c),
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

	up := new(GHReqError)
	_ = Decoder(response.Body, up)

	Trace.Println(up)

	return nil
}
