package utils

import (
	"Goonker/server/assets"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
)

var challenges []Challenge

type Challenge struct {
	Question  string   `json:"question"`
	Answers   []string `json:"answers"`
	AnswerKey int      `json:"answer_key"`
}

func LoadChallenges() {
	challengesByte, err := assets.AssetsFS.ReadFile("challenges.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(challengesByte, &challenges); err != nil {
		log.Fatal(err)
	}
}

func PickChallenge() (*Challenge, error) {
	if challenges == nil {
		return nil, fmt.Errorf("no challenges loaded")
	}

	randIndex := rand.Intn(len(challenges))
	return &challenges[randIndex], nil
}
