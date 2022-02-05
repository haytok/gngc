package gngc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dilmnqvovpnmlib/gngc/model"
	"github.com/joho/godotenv"
	"github.com/shurcooL/githubv4"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

var (
	writer io.Writer = os.Stderr
	Info             = log.New(writer, "INFO: ", log.LstdFlags)
	Error            = log.New(writer, "ERROR: ", log.LstdFlags)
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		Error.Print("Error for Loading .env file.")
		os.Exit(1)
	}
}

func PostGitHubAPI() (model.ResultContributionDays, error) {
	api_token := os.Getenv("API_TOKEN")
	user_name := os.Getenv("USER_NAME")
	if api_token == "" {
		return model.ResultContributionDays{}, errors.New("Token for GitHub API doesnt exits.")
	}
	if user_name == "" {
		return model.ResultContributionDays{}, errors.New("User name for GitHub doesnt exits.")
	}

	url := "https://api.github.com/graphql"
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: api_token})
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient(url, httpClient)

	var query struct {
		RepositoryOwner struct {
			Login graphql.String
			User  struct {
				Bio                     string
				ContributionsCollection struct {
					ContributionCalendar struct {
						TotalContributions int
						Weeks              []model.ResultWeeks
					}
				}
			} `graphql:"... on User"`
		} `graphql:"repositoryOwner(login: $userName)"`
	}
	vaiables := map[string]interface{}{
		"userName": githubv4.String(user_name),
	}

	res := model.ResultContributionDays{}
	err := client.Query(context.Background(), &query, vaiables)
	if err != nil {
		return res, err
	}

	var weeks []model.ResultWeeks
	weeks = query.RepositoryOwner.User.ContributionsCollection.ContributionCalendar.Weeks
	week := weeks[len(weeks)-1].ContributionDays
	latest_day := week[len(week)-1]
	res.ContributionCount = latest_day.ContributionCount
	res.Date = latest_day.Date

	return res, err
}

func createMessage(result model.ResultContributionDays) string {
	splitedDate := strings.Split(result.Date, "-")
	msg := splitedDate[0] + " 年 " +
		splitedDate[1] + " 月 " +
		splitedDate[2] + " 日のコミット数は " +
		strconv.Itoa(result.ContributionCount) + " です！"

	return msg
}

func GetGitHubContributions() string {
	res, err := PostGitHubAPI()
	if err != nil {
		Error.Print(err)
		os.Exit(1)
	}

	msg := createMessage(res)
	return msg
}

func createUrl(url string, event string) (string, error) {
	token := os.Getenv("IFTTT_TOKEN")
	if token == "" {
		return "", errors.New("Token for IFTTT doesnt exits.")
	}

	return fmt.Sprintf(url, event) + token, nil
}

func postIFTTT(url string, msg string) (*http.Response, error) {
	title := "GitHub Contrbution"
	jsonData := `{"value1":"` + title + `", "value2":"` + msg + `", "value3": "hoge"}`

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonData)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NotifyIFTTT(msg string) {
	url := "https://maker.ifttt.com/trigger/%s/with/key/"
	event := "tools"
	url, err := createUrl(url, event)
	if err != nil {
		Error.Print(err)
		os.Exit(1)
	}

	resp, err := postIFTTT(url, msg)
	if err != nil {
		Error.Print(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Error.Print(err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}
