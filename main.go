package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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

	r, _, err := client.Repositories.Get(ctx, owner, targetrepo)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("hit rate limit")
		os.Exit(1)
	} else if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	since := r.PushedAt.Time.Add(time.Nanosecond)

	opt := &github.RepositoryListForksOptions{
		Sort:        "newest",
		ListOptions: github.ListOptions{Page: 1, PerPage: 100},
	}

	forkedrepos := make([]*github.Repository, 0, 1000)
	for {
		repos, resp, err := client.Repositories.ListForks(ctx, owner, targetrepo, opt)
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
			os.Exit(1)
		} else if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		forkedrepos = append(forkedrepos, repos...)
		resp.Body.Close()
		log.Printf("Finish searching forked repos: %d\n", (opt.Page-1)*100)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	activerepos := make([]*github.Repository, 0, 100)
	for _, forkedrepo := range forkedrepos {
		if *forkedrepo.HasIssues {
			activerepos = append(activerepos, forkedrepo)
		}
	}

	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	writer.Write([]string{"ForkedAt", "UpdatedAt", "commits", "Star", "User", "URL"})

	for _, activerepo := range activerepos {
		if activerepo.UpdatedAt == nil {
			continue
		}
		if activerepo.Owner == nil {
			continue
		}
		if activerepo.Owner.Login == nil {
			continue
		}
		_, _, err = client.Repositories.Get(ctx, *activerepo.Owner.Login, *activerepo.Name)
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("Error(): GItHub API hit rate limit")
			os.Exit(1)
		} else if err != nil {
			continue
		}

		opt := &github.CommitsListOptions{
			Since: since,
			ListOptions: github.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		}

		log.Printf("Searcing commits: %s", *activerepo.HTMLURL)
		commits := make([]*github.RepositoryCommit, 0, 100)
		for {
			c, resp, err := client.Repositories.ListCommits(ctx, *activerepo.Owner.Login, *activerepo.Name, opt)
			if _, ok := err.(*github.RateLimitError); ok {
				log.Println("Error(): GItHub API hit rate limit")
				os.Exit(1)
			} else if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			if resp.StatusCode != 200 {
				resp.Body.Close()
				break
			}
			commits = append(commits, c...)
			log.Printf("Finish counting commits: %d\n", (opt.Page-1)*100)

			if resp.NextPage == 0 {
				resp.Body.Close()
				break
			}
			opt.Page = resp.NextPage
			resp.Body.Close()
		}

		s := []string{activerepo.CreatedAt.Format(time.RFC3339), activerepo.UpdatedAt.Format(time.RFC3339), strconv.Itoa(len(commits)), strconv.Itoa(*activerepo.StargazersCount), *activerepo.Owner.Login, *activerepo.HTMLURL}

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
		log.Println("repository format is wrong: param0 = owner/repo")
		os.Exit(1)
	}
	owner = s[0]
	repo = s[1]
	filename = os.Args[2]
	return owner, repo, filename
}
