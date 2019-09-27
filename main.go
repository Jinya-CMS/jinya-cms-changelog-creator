package main

import (
	"bytes"
	"encoding/json"
	"github.com/andlabs/ui"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"text/template"
)

var mainWindow *ui.Window
var versionsDropdown *ui.Combobox
var markdownEntry *ui.MultilineEntry
var generateChangelogButton *ui.Button

var versionBundle VersionBundle

type VersionBundleElement struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"$type"`
}

type VersionBundle struct {
	Values []VersionBundleElement `json:"values"`
	Type   string                 `json:"$type"`
}

type Issue struct {
	Summary string `json:"summary"`
	Id      string `json:"idReadable"`
	Type    string `json:"$type"`
}

type TemplateData struct {
	Issues []Issue
	Name   string
}

func loadVersions() {
	resp, err := http.Get("https://jinya.myjetbrains.com/youtrack/api/admin/customFieldSettings/bundles/version/71-0?fields=values(id,name)")
	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	err = json.Unmarshal(body, &versionBundle)
	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	for _, version := range versionBundle.Values {
		versionsDropdown.Append(version.Name)
	}
}

func generateChangeLog() {
	version := versionBundle.Values[versionsDropdown.Selected()].Name
	escapedPath := url2.PathEscape("query=project:JGCMS Fix versions:")
	url := "https://jinya.myjetbrains.com/youtrack/api/issues?fields=summary,idReadable&" + escapedPath + version
	resp, err := http.Get(url)
	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	var issues []Issue
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	err = json.Unmarshal(body, &issues)
	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	markdownTemplate :=
		`## Version {{.Name}}
### New Features

{{range .Issues}} * [{{.Id}}](https://jinya.myjetbrains.com/youtrack/issue/{{.Id}}) {{.Summary}}
{{end}}`

	templateData := TemplateData{
		Issues: issues,
		Name:   version,
	}
	tmpl, err := template.New("markdown").Parse(markdownTemplate)
	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	var markdown bytes.Buffer
	err = tmpl.Execute(&markdown, templateData)
	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	markdownEntry.SetText(markdown.String())
}

func setupUi() {
	mainWindow = ui.NewWindow("Jinya Changelog Creator", 640, 480, false)
	mainWindow.OnClosing(func(window *ui.Window) bool {
		ui.Quit()
		return true
	})

	ui.OnShouldQuit(func() bool {
		mainWindow.Destroy()
		return true
	})

	versionsDropdown = ui.NewCombobox()
	markdownEntry = ui.NewMultilineEntry()

	loadVersionsButton := ui.NewButton("Load Versions")
	loadVersionsButton.OnClicked(func(button *ui.Button) {
		go loadVersions()
	})

	versionsDropdown.OnSelected(func(combobox *ui.Combobox) {
		generateChangelogButton.Enable()
	})

	generateChangelogButton = ui.NewButton("Generate Changelog")
	generateChangelogButton.Disable()
	generateChangelogButton.OnClicked(func(button *ui.Button) {
		go generateChangeLog()
	})

	versionForm := ui.NewForm()
	versionForm.Append("Version", versionsDropdown, false)
	versionForm.SetPadded(true)

	markdownForm := ui.NewForm()
	markdownForm.Append("Changelog Markdown", markdownEntry, true)
	markdownForm.SetPadded(true)

	layoutContainer := ui.NewVerticalBox()
	layoutContainer.SetPadded(true)
	layoutContainer.Append(loadVersionsButton, false)
	layoutContainer.Append(versionForm, false)
	layoutContainer.Append(generateChangelogButton, false)
	layoutContainer.Append(markdownForm, true)

	mainWindow.SetChild(layoutContainer)
	mainWindow.SetMargined(true)

	mainWindow.Show()
}

func main() {
	ui.Main(setupUi)
}
