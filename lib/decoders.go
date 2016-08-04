package lib

import (
	"encoding/json"
	"io"
)

func Decode(r io.Reader) (users *SOUsers, err error) {

	users = new(SOUsers)
	return users, json.NewDecoder(r).Decode(users)
}

func Decode2(r io.Reader) (repo *Repo, err error) {

	repo = new(Repo)
	return repo, json.NewDecoder(r).Decode(repo)
}

func Decode3(r io.Reader) (up *GHReqError, err error) {

	up = new(GHReqError)
	return up, json.NewDecoder(r).Decode(up)
}
