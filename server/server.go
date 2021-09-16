package qserver

import (
    "net/http"
    "fmt"
)

type QuizAnswers map[string]string

type QuizQuestion struct {
    Question string `json:"question"`
    Answer QuizAnswers `json:"answer"`
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "hello world")
}
