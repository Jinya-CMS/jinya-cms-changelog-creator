package main

import (
	"github.com/andlabs/ui"
	"jinya-changelog-creator/markdown"
	"jinya-changelog-creator/youtrack"
)

var mainWindow *ui.Window
var versionsDropdown *ui.Combobox
var markdownEntry *ui.MultilineEntry
var generateChangelogButton *ui.Button

var versions []youtrack.VersionBundleElement

func generateChangeLog() {
	version := versions[versionsDropdown.Selected()].Name
	issues, err := youtrack.LoadIssues(version)

	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	templateData := markdown.TemplateData{
		Issues: issues,
		Name:   version,
	}

	result, err := markdown.ConvertMarkdown(&templateData)

	if err != nil {
		ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
		return
	}

	markdownEntry.SetText(result)
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
		versionBundle, err := youtrack.LoadVersions()

		if err != nil {
			ui.MsgBoxError(mainWindow, "Unexpected Error", err.Error())
			return
		}

		for _, version := range versionBundle.Values {
			versionsDropdown.Append(version.Name)
		}

		versions = versionBundle.Values
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
