package utils

import "github.com/AlecAivazis/survey/v2"

func AskPromptOptionList(
	message string,
	options []string,
	size int) (string, error) {

	prompt := &survey.Select{
		Message: message,
		Options: options,
	}

	answer := ""
	if err := survey.AskOne(
		prompt,
		&answer,
		survey.WithIcons(func(icons *survey.IconSet) {
			icons.SelectFocus.Format = "green+hb"
		}), survey.WithPageSize(size)); err != nil {
		return "No", err
	}

	return answer, nil
}

func AskYesOrNo(Message string) (string, error) {
	prompt := &survey.Select{
		Message: Message,
		Options: []string{"Yes", "No (exit)"},
	}

	answer := ""
	if err := survey.AskOne(prompt, &answer, survey.WithIcons(func(icons *survey.IconSet) {
		icons.SelectFocus.Format = "green+hb"
	}), survey.WithPageSize(2)); err != nil {
		return "No", err
	}

	return answer, nil
}
