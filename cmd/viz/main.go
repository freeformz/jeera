package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/freeformz/jeera/cmd/internal/config"
)

func jiraClient() (*jira.Client, error) {
	cfg, token, err := config.OAuth()
	if err != nil {
		return nil, err
	}
	url, err := config.JiraURL()
	if err != nil {
		return nil, err
	}

	return jira.NewClient(
		cfg.Client(context.Background(), token),
		url,
	)
}

func projectForIssue(key string) (string, error) {
	p := strings.SplitN(key, "-", 2)
	if len(p) < 2 {
		return "", fmt.Errorf("Malformed issue (%s), expected <PROJECT>-<NUMBER>", key)
	}
	return p[0], nil
}

func main() {
	client, err := jiraClient()
	if err != nil {
		log.Fatal("setting up jira client:", err)
	}
	if len(os.Args) < 2 {
		log.Fatal("first arg (issue key) can't be blank")
	}

	ek := os.Args[1]
	p, err := projectForIssue(ek)
	if err != nil {
		log.Fatal(err)
	}
	projects := make(map[string][]*jira.Issue)

	cfID, err := config.JiraCustomEpicFieldID()
	if err != nil {
		log.Fatal("determining epic custom field id:", err)
	}

	e, _, err := client.Issue.Get(ek, nil)
	if err != nil {
		log.Fatal("fetching issue:", err)
	}

	if e.Fields.Type.Name != "Epic" {
		log.Fatal("Issue " + ek + "is not an epic")
	}

	projects[p] = []*jira.Issue{e}

	is, _, err := client.Issue.Search("cf["+cfID+"]="+ek, nil)
	if err != nil {
		log.Fatal("searching for epic stories:", err)
	}
	for _, i := range is {
		p, err := projectForIssue(i.Key)
		if err != nil {
			log.Fatalf("determining project for issue %s in epic %s: %s", i.Key, ek, err.Error())
		}
		issues, ok := projects[p]
		if !ok {
			issues = []*jira.Issue{&i}
			projects[p] = issues
		}
		fmt.Println(i.Key, i.Fields.Status.Name)
	}

	for _, l := range e.Fields.IssueLinks {
		if l.OutwardIssue != nil {
			fmt.Println("->", l.OutwardIssue.Key)
		}
		if l.InwardIssue != nil {
			fmt.Println("<-", l.InwardIssue.Key)
		}
	}

	fmt.Println("graph {")

}
