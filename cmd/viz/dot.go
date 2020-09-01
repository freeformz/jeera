package main

import (
	"path"

	"github.com/andygrunwald/go-jira"
	"github.com/emicklei/dot"
)

func addToGraph(g *dot.Graph, i *jira.Issue, url string) (dot.Node, error) {
	p, err := projectForIssue(i.Key)
	if err != nil {
		return dot.Node{}, err
	}
	sg := g.Subgraph(p, dot.ClusterOption{})
	n := sg.Node(i.Key).
		Attr("URL", path.Join(url, "browse", i.Key)).
		Attr("tooltip", i.Fields.Summary)

	//log.Println("Type: ", i.Fields.Type.Name)
	switch i.Fields.Type.Name {
	case "Epic":
		n = n.Attr("shape", "doubleoctagon")
	default:
		n = n.Attr("shape", "circle")
	}

	switch i.Fields.Status.StatusCategory.Name {
	case "Done":
		n = n.Attr("fillcolor", "green").
			Attr("style", "filled")
	case "In Progress":
		n = n.Attr("fillcolor", "yellow").
			Attr("style", "filled")
	case "To Do":
		n = n.Attr("fillcolor", "orange").
			Attr("style", "filled")
	default:
	}

	for _, l := range i.Fields.IssueLinks {
		if l.OutwardIssue != nil {
			on, err := addToGraph(g, l.OutwardIssue, url)
			if err != nil {
				return n, err
			}
			if e := g.FindEdges(n, on); len(e) == 0 {
				decorateEdge(n.Edge(on), l)
			}
		}
		if l.InwardIssue != nil {
			in, err := addToGraph(g, l.InwardIssue, url)
			if err != nil {
				return n, err
			}
			if e := g.FindEdges(in, n); len(e) == 0 {
				decorateEdge(in.Edge(n), l)
			}
		}
	}

	return n, nil
}

func decorateEdge(e dot.Edge, l *jira.IssueLink) {
	e.Label(l.Type.Name)
	switch l.Type.Name {
	case "Blocks":
		e.Attr("color", "red")
	}
}
