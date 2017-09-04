package main

import (
	"golang.org/x/oauth2"
	"github.com/google/go-github/github"
	"context"
	"fmt"
	"encoding/json"
	"regexp"
	"sort"
	"strings"
	time2 "time"
	"os"
)

const errorMessage string = "json.Unmarshal failed with '%s'\n"

var(
	yourName = "Tomohiro Imaizumi"
)

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {
	var personalAccessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	if personalAccessToken == "" {
		fmt.Println("Set access token as an envvar of \"GITHUB_ACCESS_TOKEN\".")
		return
	}
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := github.NewClient(oauthClient)
	DoFetchEvents(client)
}

func DoFetchEvents(client *github.Client) {
	var orgName = os.Getenv("ORG_NAME")
	if orgName == "" {
		fmt.Println("Set your organization name as an envvar of \"ORG_NAME\".")
		return
	}
	var repoName = os.Getenv("REPO_NAME")
	if repoName == "" {
		fmt.Println("Set your repogitory name as an envvar of \"REPO_NAME\".")
		return
	}

	pushes := []github.Event{}
	comments := []github.Event{}
	reviews := []github.Event{}
	for i := 0; ; {
		const pageMax int = 2
		if i > pageMax {
			break
		}
		activity := client.Activity
		evt, _, clientError := activity.ListRepositoryEvents(
			context.Background(),
			orgName,
			repoName,
			&github.ListOptions{Page: i + 1, PerPage: 100})
		i++
		if clientError != nil {
			fmt.Printf("client.Users.Get() faled with '%s'\n", clientError)
			continue
		}
		pushEvents := Filter(evt, "PushEvent", IsType)
		pushes = append(pushes, pushEvents...)

		issueCommentEvents := Filter(evt, "IssueCommentEvent", IsType)
		comments = append(comments, issueCommentEvents...)

		pullRequestReviewCommentEvents := Filter(evt, "PullRequestReviewCommentEvent", IsType)
		reviews = append(reviews, pullRequestReviewCommentEvents...)
	}
  fmt.Printf("\n# Daily Report: %s %s\n", yourName, time2.Now().Format("2006-01-02"))
	fmt.Printf("\n## Commits\n\n")
	ParseCommits(&pushes)
	fmt.Printf("\n")
	fmt.Printf("\n## Comments\n\n")
	ParseComments(&comments)
}

func ParseComments(comments *[]github.Event) {
	for _, comment := range *comments {
		if comment.Actor.GetLogin() != "imaizume" { continue }
		creationTime := comment.CreatedAt.Local()
		limitTime := time2.Now().Add(-time2.Hour * 18)
		if creationTime.Before(limitTime) { continue }
		var t github.IssueCommentEvent
		marshalError := json.Unmarshal(*comment.RawPayload, &t)
		if marshalError != nil {
			fmt.Printf(errorMessage, marshalError)
			return
		}
		comment := t.Comment.GetBody()
		line := strings.Replace(comment, "\n", " ", -1)
		line2 := strings.Replace(line, "\r\n", " ", -1)
		line3 := strings.Replace(line2, "\r", " ", -1)
		creation := t.Comment.CreatedAt.Local()
		fmt.Printf("- [%d:%d] %s\n", creation.Hour(), creation.Minute(), line3)
		// fmt.Printf("[%d:%d] %s", creation.Hour(), creation.Minute(), line)
	}
	fmt.Printf("\n### Note\n\n")
}

func ParseCommits(pushes *[]github.Event) {
	var evt GHEvent = *pushes
	sort.Sort(evt)
	var lastBranchName string = ""
	for _, push := range *pushes {
		if push.Actor.GetLogin() != "imaizume" { continue }
		creationTime := push.CreatedAt.Local()
		limitTime := time2.Now().Add(-time2.Hour * 18)
		if creationTime.Before(limitTime) { continue }
		var t github.PushEvent
		marshalError := json.Unmarshal(*push.RawPayload, &t)
		if marshalError != nil {
			fmt.Printf(errorMessage, marshalError)
			return
		}
		branchPattern := regexp.MustCompile(`(imaizumi\.\d+\.[\d\w_\.]+|develop\.v4\.0\.2\.master\.merge)$`)
		branchName := t.GetRef()
		matches := branchPattern.FindAllStringSubmatch(branchName, -1)
		if len(matches) > 0 {
			if branchName != lastBranchName {
				fmt.Printf("### %s\n\n", matches[0][0])
			}
			lastBranchName = branchName
			var lines = DissolvePushToCommits(&t, yourName)
			for _, l := range lines {
				fmt.Println(l)
			}
		}
		fmt.Printf("\n### Note\n\n")
	}
}

func DissolvePushToCommits(push *github.PushEvent, username string) []string {
	var cmt GHPushEventCommit = push.Commits
	sort.Sort(cmt)
	lines := []string{}
	for _, v := range cmt {
		message := strings.Split(v.GetMessage(), "\n")
		if v.Author.GetName() != username { continue }
		lines = append(lines, fmt.Sprintf("- (%s) %s", v.GetSHA()[0:6], message[0]))
	}
	return lines
}

func IsType(event github.Event, t string) bool {
	return event.GetType() == t
}

func Filter(vs []*github.Event, t string, f func(github.Event, string) bool) []github.Event {
	vsf := make([]github.Event, 0)
	for _, v := range vs {
		if f(*v, t) {
			vsf = append(vsf, *v)
		}
	}
	return vsf
}

type Psh []github.PushEvent
type GHPushEventCommit []github.PushEventCommit

func (e GHPushEventCommit) Len() int {
	return len(e)
}

func (p GHPushEventCommit) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (event GHPushEventCommit) Less(i, j int) bool {
	return event[i].GetTimestamp().Time.After(event[i].GetTimestamp().Time)
}
