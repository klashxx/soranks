package lib

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func GitHubConnector(md string, branch string, author Committer) error {

	encoded, err := Markdown2Base64(RspMDPath)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", GHApiURL, "/git/trees/dev")
	Trace.Printf("Tree url: %s\n", url)

	folder := false
	repo := new(Repo)

	err = StreamHTTP(url, repo, false)
	if err != nil {
		return err
	}

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

	err = StreamHTTP(url, repo, false)
	if err != nil {
		return err
	}

	sha := ""
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
		return fmt.Errorf("Can't get github  token!")
	}
	Info.Printf("token: %s\n", token)

	c := fmt.Sprintf("%s [%s]", md, time.Now().Format(time.RFC3339))

	var buf io.ReadWriter

	if sha == "" {
		Info.Println("Update not detected.")
		data := Create{
			Path:      RspMDPath,
			Message:   fmt.Sprintf("Create: %s", c),
			Content:   encoded,
			Branch:    branch,
			Committer: author}
		buf, err = Encoder(data)
	} else {
		Info.Printf("Update SHA: %s", sha)
		data := Update{
			Path:      RspMDPath,
			Message:   fmt.Sprintf("Update: %s", c),
			Content:   encoded,
			Sha:       sha,
			Branch:    branch,
			Committer: author}
		buf, err = Encoder(data)
	}

	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, buf)
	if err != nil {
		return err
	}
	Trace.Println("Sending header.")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	Trace.Println("Response.")
	defer response.Body.Close()

	up := new(GHReqError)
	err = Decoder(response.Body, up)
	if err != nil {
		return err
	}

	Trace.Println(up)

	return nil
}
