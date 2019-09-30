package markdown

import "jinya-changelog-creator/youtrack"

type TemplateData struct {
	Type   string
	Issues []youtrack.Issue
}
