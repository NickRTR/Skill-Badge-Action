package badge

import "fmt"

func Format(skills []SkillSection) string {
	var markdown string

	for _, skillSection := range skills {
		markdown += fmt.Sprintf("\n### %s\n\n", skillSection.Title)
		for _, skill := range skillSection.Skills {
			markdown += fmt.Sprintf("<img alt=\"%s\" src=\"%s\"/> ", skill.Title, skill.URL)
		}
	}

	return markdown
}
