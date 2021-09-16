package request

type GetQuestionRequest struct {
    User string `json:"user"`
}

type AnswerRequest struct {
    User string `json:"user"`
    Answer string `json:"answer"`
}
