package ui

import (
	"github.com/charmbracelet/huh"
)

func PromptTodoLine(title string, value string) (string, error) {
	var input string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title(title).Prompt("> ").Placeholder("todo.txt task").Value(&input),
		),
	)
	if value != "" {
		input = value
	}
	if err := form.Run(); err != nil {
		return "", err
	}
	return input, nil
}
