package qserver

import (
    "testing"
    "net/http/httptest"
    "io/ioutil"
    "fmt"
    "strings"
    "encoding/json"
)

func TestGetQuestion(t *testing.T) {
    wanted, _ := json.Marshal(QuizQuestion{
        Question: "some question",
        Answer: QuizAnswers{
            "1": "one",
            "2": "two",
            "3": "three",
        },
    })

    requestBody := strings.NewReader(`{"user":"1234"}`)
    request := httptest.NewRequest("POST", "http://127.0.0.1/get_question", requestBody)
    recorder := httptest.NewRecorder()
    GetQuestion(recorder, request)


    recorder.Write(wanted)
    response := recorder.Result() 
    result, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println(err)
    }

    if string(result) != string(wanted) {
        t.Fatalf("\ngot: %q\nwanted:%q", string(result), string(wanted))
    }
}
