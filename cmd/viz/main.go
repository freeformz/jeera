package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/emicklei/dot"
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
	url, err := config.JiraURL()
	if err != nil {
		log.Fatal(err)
	}
	client, err := jiraClient()
	if err != nil {
		log.Fatal("setting up jira client:", err)
	}
	if len(os.Args) < 2 {
		log.Fatal("first arg (issue key) can't be blank")
	}

	ek := os.Args[1]
	if _, err := projectForIssue(ek); err != nil {
		log.Fatal(err)
	}

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
	eg := dot.NewGraph(dot.Directed).
		Label(ek + " Graph")
	eg.Attr("rankdir", "LR")
	eg.Attr("labelloc", "t")
	//eg.Attr("clusterrank", "local")

	is, _, err := client.Issue.Search("cf["+cfID+"]="+ek, nil)
	if err != nil {
		log.Fatal("searching for epic stories:", err)
	}
	for _, i := range is {
		_, err := addToGraph(eg, &i, url)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(eg.String())
}
