package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/machinebox/graphql"
)

var client *graphql.Client

func setupClient(URL string) {
	client = graphql.NewClient(URL)
}

func getAllSkills(TOKEN string) {
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

	var responseData SkillSectionsResponse

	if err := client.Run(ctx, query, &responseData); err != nil {
		log.Fatal(err)
	}

	for _, skillSection := range responseData.SkillSections {
		fmt.Println(skillSection.Title)
	}
}

func init() {
	// environment variables for local development
	if len(os.Args) > 1 {
		if os.Args[1] == "dev" {
			error := godotenv.Load(".env")
			if error != nil {
				log.Fatalln("Could not load .env file")
			}
		}
	}
}

func main() {
	var URL string = os.Getenv("HYGRAPH_API_URL")
	var TOKEN string = os.Getenv("HYGRAPH_API_TOKEN")
	// var GH_TOKEN string = os.Getenv("GH_TOKEN")
	// target := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	setupClient(URL)
	getAllSkills(TOKEN)
}
