package api

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetGithubUser(c *gin.Context, config *oauth2.Config, code string) (user *github.User, err error) {
	token, err := config.Exchange(c, code)
	if err != nil {
		return
	}
	client := config.Client(context.TODO(), token)
	response, err := client.Get("https://api.github.com/user")
	if err != nil {
		return
	}
	defer response.Body.Close()
	info, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(info, &user)
	if err != nil {
		return
	}
	return
}
