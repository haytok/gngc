package gngc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/dilmnqvovpnmlib/gngc/model"
	"github.com/shurcooL/githubv4"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

func postGitHubAPI(gitHubConfig model.GitHubConfig) (model.ResultContributionDays, error) {
	url := "https://api.github.com/graphql"
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gitHubConfig.Token})
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
		"userName": githubv4.String(gitHubConfig.UserName),
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

func GetGitHubContributions(gitHubConfig model.GitHubConfig) (string, error) {
	res, err := postGitHubAPI(gitHubConfig)
	if err != nil {
		return "", err
	}

	msg := createMessage(res)

	return msg, nil
}

func createUrl(iFTTTConfig model.IFTTTConfig, url string) string {
	return fmt.Sprintf(url, iFTTTConfig.EventName) + iFTTTConfig.Token
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

func NotifyIFTTT(iFTTTConfig model.IFTTTConfig, msg string) (string, error) {
	url := "https://maker.ifttt.com/trigger/%s/with/key/"
	url = createUrl(iFTTTConfig, url)

	resp, err := postIFTTT(url, msg)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res_message := resp.Status + " from IFTTT API\n" + string(body)

	return res_message, nil
}
