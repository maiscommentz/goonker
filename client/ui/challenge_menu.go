package ui

import (
	"Goonker/common"
)

const (
	AnswerButtonY        = 150
	ButtonSpacingY       = 20
	ChallengeButtonWidth = 400
)

// ChallengeMenu represents the UI for a challenge.
type ChallengeMenu struct {
	Question string
	Answers  []Button
	Clock    Timer
}

// NewChallengeMenu creates a new ChallengeMenu instance.
func NewChallengeMenu(challenge common.ChallengePayload) *ChallengeMenu {
	challengeMenu := &ChallengeMenu{Question: challenge.Question}

	// Center buttons
	centerX := (float64(WindowWidth) - ChallengeButtonWidth) / 2

	// Answer buttons
	for i, answer := range challenge.Answers {
		buttonHeight := float64(AnswerButtonY + i*(ButtonHeight+ButtonSpacingY))
		btn := NewButton(centerX, buttonHeight, ChallengeButtonWidth, ButtonHeight, answer, SmallFontFace)
		challengeMenu.Answers = append(challengeMenu.Answers, *btn)
	}

	return challengeMenu
}
