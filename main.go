package main

import (
	"github.com/andlabs/ui"
	"jinya-changelog-creator/markdown"
	"jinya-changelog-creator/youtrack"
	"sort"
)

var mainWindow *ui.Window
var versionsDropdown *ui.Combobox
var markdownEntry *ui.MultilineEntry
var generateChangelogButton *ui.Button

var versions []string

func loadData() {
	versionBundle, err := youtrack.LoadVersions()

	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	versions = make([]string, len(versionBundle.Values))

	for idx, version := range versionBundle.Values {
		versions[idx] = version.Name
	}

	sort.Strings(versions)

	for _, version := range versions {
		versionsDropdown.Append(version)
	}
}

func generateChangeLog() {
	markdownEntry.SetText("")
	version := versions[versionsDropdown.Selected()]
	issueTypes, err := youtrack.LoadIssueTypes()

	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	markdownEntry.Append("## Version " + version)
	markdownEntry.Append("\n")

	for _, issueType := range issueTypes {
		issues, err := youtrack.LoadIssues(version, issueType)

		if err != nil {
			ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
			return
		}

		if len(issues) > 0 {
			templateData := markdown.TemplateData{
				Issues: issues,
				Type:   issueType,
			}

			result, err := markdown.ConvertMarkdown(&templateData)
			if err != nil {
				ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
				return
			}

			markdownEntry.Append(result)
		}
	}
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
		go loadData()
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
	err := ui.Main(setupUi)

	if err != nil {
		panic(err)
	}
}
