package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/NickRTR/Skill-Badge-Action/badge"
	"github.com/NickRTR/Skill-Badge-Action/cli"
	"github.com/NickRTR/Skill-Badge-Action/github"
	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
)

var client *graphql.Client

func setupClient(URL string) {
	client = graphql.NewClient(URL)
}

func updatedSinceLastWeek(skills []badge.SkillSection) bool {
	for _, skillSection := range skills {
		date, err := time.Parse(time.RFC3339Nano, skillSection.UpdatedAt)
		if err != nil {
			cli.BrintErr(err.Error())
			os.Exit(1)
		}
		today := time.Now()
		daysDiff := today.Sub(date).Hours() / 24

		if daysDiff <= 7 {
			return true
		}
	}
	return false
}

func getAllSkills(TOKEN string) []badge.SkillSection {
	query := graphql.NewRequest(`
		query {
			skillSections {
				title
				skills {
					title
					url
				}
				updatedAt
			}
		}
	`)

	query.Header.Set("Authorization", "Bearer "+TOKEN)

	ctx := context.Background()

	var responseData badge.SkillSectionsResponse

	if err := client.Run(ctx, query, &responseData); err != nil {
		cli.BrintErr(err.Error())
		os.Exit(1)
	}

	return responseData.SkillSections
}

func init() {
	// environment variables for local development
	if len(os.Args) > 1 {
		if os.Args[1] == "dev" {
			error := godotenv.Load(".env")
			if error != nil {
				cli.BrintErr("Could not load .env file")
			}
		}
	}
}

func main() {
	var URL string = os.Getenv("HYGRAPH_API_URL")
	var TOKEN string = os.Getenv("HYGRAPH_API_TOKEN")
	var GH_TOKEN string = os.Getenv("GH_TOKEN")
	target := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")

	setupClient(URL)
	skills := getAllSkills(TOKEN)

	if !updatedSinceLastWeek(skills) {
		cli.Brint("Skills haven't been updated since last run -> Action done.")
		os.Exit(0)
	}

	markdown := badge.Format(skills)

	client := github.Authenticate(GH_TOKEN)
	github.AddBadges(client, markdown, target[0], target[1])
}
