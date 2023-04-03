package main

type Skill struct {
	Title string
	URL   string
}

type SkillSections struct {
	Title  string
	Skills []Skill
}

type SkillSectionsResponse struct {
	SkillSections []SkillSections
}
