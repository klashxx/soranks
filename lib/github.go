package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func GetToken() (token string) {
	Info.Println("Trying to load TOKEN.")
	token = os.Getenv("GH_TOKEN")
	if token != "" {
		return token
	}
	Warning.Println("Can't load GH_TOKEN env variable")

	Info.Println("Trying to extract TOKEN.")
	token, err := GetKey(gitHubToken)
	if err != nil {
		Error.Println(err)
		os.Exit(5)
	}
	return token
}

func GHPublisher(token string, publish *string, branch string, author Committer) error {
	fname := fmt.Sprintf("%s.md", *publish)
	if err := GitHubConnector(rspMDpath, fname, token, branch, author); err != nil {
		return fmt.Errorf("GitHub connection Markdown (%s) error: %s\n", fname, err)
	}

	fname = fmt.Sprintf("%s.json", *publish)
	if err := GitHubConnector(rspJSONpath, fname, token, branch, author); err != nil {
		return fmt.Errorf("GitHub connection JSON (%s) error: %s\n", fname, err)
	}
	return nil
}

func GitHubConnector(fmtpath string, target string, token string, branch string, author Committer) error {

	encoded, err := F2Base64(fmtpath)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s%s", ghAPIURL, "/git/trees/dev")
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

	Trace.Printf("target: %s\n", target)

	err = StreamHTTP(url, repo, false)
	if err != nil {
		return err
	}

	sha := ""
	for _, file := range repo.Tree {
		if file.Path == target {
			sha = file.Sha
			break
		}
	}

	url = fmt.Sprintf("%s/contents/data/%s", ghAPIURL, target)
	Trace.Println(url)

	c := fmt.Sprintf("%s [%s]", target, time.Now().Format(time.RFC3339))

	var buf io.ReadWriter

	if sha == "" {
		Info.Println("Update not detected.")
		data := Create{
			Path:      fmtpath,
			Message:   fmt.Sprintf("Create: %s", c),
			Content:   encoded,
			Branch:    branch,
			Committer: author}
		buf, err = JSONEncoder(data)
	} else {
		Info.Printf("Update SHA: %s", sha)
		data := Update{
			Path:      fmtpath,
			Message:   fmt.Sprintf("Update: %s", c),
			Content:   encoded,
			Sha:       sha,
			Branch:    branch,
			Committer: author}
		buf, err = JSONEncoder(data)
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
	err = JSONDecoder(response.Body, up)
	if err != nil {
		return err
	}

	Trace.Println(up)

	return nil
}
