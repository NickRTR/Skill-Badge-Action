package badge

type Skill struct {
	Title string
	URL   string
}

type SkillSection struct {
	Title  string
	Skills []Skill
}

type SkillSectionsResponse struct {
	SkillSections []SkillSection
}
