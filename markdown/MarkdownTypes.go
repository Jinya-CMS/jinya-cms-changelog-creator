package markdown

import "jinya-changelog-creator/youtrack"

type TemplateData struct {
	Issues []youtrack.Issue
	Name   string
}
