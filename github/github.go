package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/NickRTR/Skill-Badge-Action/cli"
	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
)

func Authenticate(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return client
}

func AddBadges(client *github.Client, markdown string, user string, repo string) {
	file, _, _, err := client.Repositories.GetContents(context.Background(), user, repo, "README.md", nil)
	if err != nil {
		cli.BrintErr(err.Error())
	}

	readme, err := file.GetContent()
	if err != nil {
		cli.BrintErr(err.Error())
	}

	start := "<!--Skill-Badge-Action-Start-->"
	stop := "<!--Skill-Badge-Action-End-->"
	startIndex := strings.Index(readme, start)
	stopIndex := strings.Index(readme, stop) + len(stop)
	sectionBefore := readme[0:startIndex]
	sectionAfter := readme[stopIndex:]

	editedReadme := sectionBefore + fmt.Sprintf("%s\n%s\n%s", start, markdown, stop) + sectionAfter

	b := []byte(editedReadme)
	sha := file.GetSHA()
	message := "Update README skills section"

	updatedFile := github.RepositoryContentFileOptions{
		Message: &message,
		Content: b,
		SHA:     &sha,
	}

	_, _, err = client.Repositories.UpdateFile(context.Background(), user, repo, "README.md", &updatedFile)
	if err != nil {
		cli.BrintErr(err.Error())
	}

	cli.Brint("Updated README skills")
}
