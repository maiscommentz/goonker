package ui

import (
	"Goonker/common"
)

const (
	AnswerButtonY  = 150
	ButtonSpacingY = 20
)

type ChallengeMenu struct {
	Question string
	Answers  []Button
	//Duration *time.Duration
	Clock Timer
}

func NewChallengeMenu(challenge common.ChallengePayload) *ChallengeMenu {
	challengeMenu := &ChallengeMenu{Question: challenge.Question}

	// Center buttons
	centerX := (float64(WindowWidth) - ButtonWidth) / 2

	// Answer buttons
	for i, answer := range challenge.Answers {
		buttonHeight := float64(AnswerButtonY + i*(ButtonHeight+ButtonSpacingY))
		btn := NewButton(centerX, buttonHeight, ButtonWidth, ButtonHeight, answer, SubtitleFontSize)
		challengeMenu.Answers = append(challengeMenu.Answers, *btn)
	}

	return challengeMenu
}
