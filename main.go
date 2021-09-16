package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/quiztest/quiz/config"
	errorModel "github.com/quiztest/quiz/models/error"
	requestModel "github.com/quiztest/quiz/models/request"
	responseModel "github.com/quiztest/quiz/models/response"
	"github.com/quiztest/quiz/storage"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var configFilePath = flag.String("config", "config.json", "path to configuration file")
var database *storage.Storage

func init() {
	flag.Parse()
	conf := config.GetConfigFromFile(*configFilePath)
	postgres := &storage.Postgres{conf}
	database = storage.CreateStorage(postgres)
}

func main() {

	http.HandleFunc("/get_question", getQuestion)
	http.HandleFunc("/answer", answer)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

var getQuestion = func(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		MethodError(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		InternalServerError(w, err)
		return
	}
	req := requestModel.GetQuestionRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Print(err)
		BadRequestError(w)
		return
	}

	user, err := database.GetUser(req.User)
	if err != nil {
		log.Print(err)
		InternalServerError(w, err)
		return
	}
	if user.Id == 0 {
		log.Printf("User %v not found", user.Id)
		NotFoundError(w)
		return
	}

	question, err := database.GetRandomQuestion()
	if err != nil {
		log.Print(err)
		InternalServerError(w, err)
		return
	}

	rand.Shuffle(len(question.Answer), func(i, j int) {
		question.Answer[i], question.Answer[j] = question.Answer[j], question.Answer[i]
	})

	responseBody, err := json.MarshalIndent(question, "", "\t")
	if err != nil {
		log.Print(err)
		InternalServerError(w, err)
		return
	}

	err = database.SaveGame(question, user)
	if err != nil {
		log.Print(err)
		InternalServerError(w, err)
		return
	}

	fmt.Fprintf(w, string(responseBody))
}

var answer = func(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		MethodError(w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		InternalServerError(w, err)
		return
	}
	req := requestModel.AnswerRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Print(err)
		BadRequestError(w)
		return
	}

	user, err := database.GetUser(req.User)
	if err != nil {
		log.Print(err)
		InternalServerError(w, err)
		return
	}
	if user.Id == 0 {
		log.Printf("User %v not found", user.Id)
		NotFoundError(w)
		return
	}

	rightAnswer, err := database.CheckAnswer(user.Id)
	if err != nil {
		log.Print(err)
		InternalServerError(w, err)
		return
	}

	isRight := strconv.Itoa(rightAnswer) == req.Answer
	response := responseModel.AnswerResponse{Right: isRight}

	responseBody, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		InternalServerError(w, err)
		return
	}

	err = database.MarkAnswered(user.Id)
	if err != nil {
		log.Print(err)
		InternalServerError(w, err)
		return
	}

	fmt.Fprintf(w, string(responseBody))
}

func InternalServerError(w http.ResponseWriter, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Internal Server Error")
}

func MethodError(w http.ResponseWriter) {
	CommonError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	return
}

func BadRequestError(w http.ResponseWriter) {
	CommonError(w, http.StatusBadRequest, "Bad Request")
	return
}

func NotFoundError(w http.ResponseWriter) {
	CommonError(w, http.StatusNotFound, "Not Found")
	return
}

func CommonError(w http.ResponseWriter, code int, message string) {
	e := errorModel.Error{
		ErrorCode:    code,
		ErrorMessage: message,
	}
	w.WriteHeader(code)
	responseBody, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		InternalServerError(w, err)
		return
	}
	fmt.Fprintf(w, string(responseBody))
	return
}
