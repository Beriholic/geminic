package ui

import (
	"github.com/beriholic/geminic/internal/service"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

type action string

var base *huh.Theme = huh.ThemeBase()

const (
	CONFIRM     action = "CONFIRM"
	REGENERATE  action = "REGENERATE"
	EDIT_COMMIT action = "EDIT_COMMIT"
	CANCEL      action = "CANCEL"
)

func RenderActionForm() (action, error) {
	var curAction action

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[action]().
				Title("Is that what you want?").
				Options(
					huh.NewOption("Yes", CONFIRM),
					huh.NewOption("Roll", REGENERATE),
					huh.NewOption("Edit", EDIT_COMMIT),
					huh.NewOption("No", CANCEL),
				).
				Value(&curAction).
				WithTheme(base),
		))

	if err := form.Run(); err != nil {
		return CANCEL, err
	}

	return curAction, nil
}

func RenderEditorForm(commit string) (action, error) {
	var confirmEdit bool = false

	input := huh.NewForm(
		huh.NewGroup(
			huh.NewText().Title("Edit commit message").CharLimit(200).Value(&commit),
		),
	)

	confirm := huh.NewConfirm().
		Title("Confirm edit?").
		Affirmative("Yes").
		Negative("No").
		Value(&confirmEdit).
		WithTheme(base)

	if err := input.Run(); err != nil {
		return CANCEL, err
	}

	if err := confirm.Run(); err != nil {
		return CANCEL, err
	}

	if !confirmEdit {
		return CANCEL, nil
	}

	return CONFIRM, service.GetGitService().CommitChanges(commit)
}

func RenderSpinner(title string, action func()) error {
	return spinner.New().
		Title(title).
		Action(action).
		Run()
}
