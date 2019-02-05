package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Github API token help: https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
// go run github-repo-export.go -h
// go run github-repo-export.go -t=qwertyuioasdfghjkl1234567890 -o=applift
func main() {
	githubToken := flag.String("t", "<YOUR GITHUB PERSONAL TOKEN>", "Your github personal token")
	organization := flag.String("o", "applift", "Your organization on Github")
	flag.Parse()
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, *organization, opt)
		if err != nil {
			panic(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	fmt.Println("repo,url,is private,is fork,is archived,stars,forks,language,created at,updated at")
	for _, repo := range allRepos {
		fmt.Printf("%s,%s,%t,%t,%t,%d,%d,%s,%s,%s\n", repo.GetName(), repo.GetURL(), repo.GetPrivate(), repo.GetArchived(), repo.GetFork(), repo.GetStargazersCount(), repo.GetForksCount(), repo.GetLanguage(), repo.GetCreatedAt(), repo.GetUpdatedAt())
	}
}
