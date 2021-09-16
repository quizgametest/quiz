package question

import (
	"encoding/json"
	"strconv"
)

type QuestionObject struct {
	Id       int
	Question string
	Answer   []AnswerObject
}

type AnswerObject struct {
	QuestionId int
	Answer     string
	IsCorrect  bool
}

func (q QuestionObject) MarshalJSON() ([]byte, error) {
	responseQuestion := struct {
		Question string            `json:"question"`
		Answer   map[string]string `json:"answer"`
	}{Question: q.Question, Answer: map[string]string{}}

	for i, a := range q.Answer {
		responseQuestion.Answer[strconv.Itoa(i+1)] = a.Answer
	}

	return json.Marshal(responseQuestion)
}
