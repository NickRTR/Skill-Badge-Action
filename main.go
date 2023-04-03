package main

import (
	"context"
	"log"
	"os"
	"strings"

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

func getAllSkills(TOKEN string) []badge.SkillSection {
	query := graphql.NewRequest(`
		query {
			skillSections {
				title
				skills {
					title
					url
				}
			}
		}
	`)

	query.Header.Set("Authorization", "Bearer "+TOKEN)

	ctx := context.Background()

	var responseData badge.SkillSectionsResponse

	if err := client.Run(ctx, query, &responseData); err != nil {
		log.Fatal(err)
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

	markdown := badge.Format(skills)

	client := github.Authenticate(GH_TOKEN)
	github.AddBadges(client, markdown, target[0], target[1])
}
