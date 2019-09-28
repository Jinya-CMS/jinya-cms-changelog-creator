package markdown

import (
	"bytes"
	"text/template"
)

func ConvertMarkdown(templateData *TemplateData) (string, error) {
	markdownTemplate :=
		`## Version {{.Name}}
### New Features

{{range .Issues}} * [{{.Id}}](https://jinya.myjetbrains.com/youtrack/issue/{{.Id}}) {{.Summary}}
{{end}}`

	tmpl, err := template.New("markdown").Parse(markdownTemplate)
	if err != nil {
		return "", err
	}

	var markdown bytes.Buffer
	err = tmpl.Execute(&markdown, templateData)
	if err != nil {
		return "", err
	}

	return markdown.String(), nil
}
