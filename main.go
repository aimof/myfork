package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	owner, targetrepo, filename := args()

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Cannot open %s\n", filename)
		os.Exit(1)
	}
	defer file.Close()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("MYFORK_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	forkedrepos, _, err := client.Repositories.ListForks(ctx, owner, targetrepo, nil)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("hit rate limit")
		os.Exit(1)
	} else if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	activerepos := make([]*github.Repository, 0, len(forkedrepos)/4)
	for _, forkedrepo := range forkedrepos {
		if *forkedrepo.HasIssues {
			activerepos = append(activerepos, forkedrepo)
		}
	}

	writer := csv.NewWriter(file)

	for _, activerepo := range activerepos {
		s := []string{activerepo.UpdatedAt.String(), strconv.Itoa(*activerepo.StargazersCount), *activerepo.URL}
		writer.Comma = '\t'
		writer.Write(s)
	}
	writer.Flush()
}

func args() (owner, repo, filename string) {
	if len(os.Args) != 3 {
		log.Println("Usage: myfork owner/repo outputfile")
		os.Exit(1)
	}
	s := strings.Split(os.Args[1], "/")
	if len(s) != 2 {
		log.Println("repository format is wrong: args1 = owner/repo")
		os.Exit(1)
	}
	owner = s[0]
	repo = s[1]
	filename = os.Args[2]
	return owner, repo, filename
}
