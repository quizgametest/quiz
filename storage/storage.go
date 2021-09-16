package storage

import (
    "fmt"
	"database/sql"
	_ "github.com/lib/pq"
    "github.com/quiztest/quiz/config"
    "github.com/quiztest/quiz/models/question"
    "github.com/quiztest/quiz/models/user"
)

type Storage struct {
    client Client
}

type Client interface {
    getConnection() (*sql.DB, error)
}

type Postgres struct {
    Config *config.Config    
}

func (p *Postgres) getConnection() (*sql.DB, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.Config.Username,
		p.Config.Password,
		p.Config.Host,
		p.Config.Port,
		p.Config.Database)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
        return nil, err
	}

    return db, nil
}

func CreateStorage(c Client) *Storage {
    return &Storage{c}
}

func (s *Storage) GetRandomQuestion() (question.QuestionObject, error) {
    db, err := s.client.getConnection()
    defer db.Close()
	if err != nil {
        return question.QuestionObject{}, err
	}
    query := `select id, question from questions offset random() * (select count(*) from questions) limit 1`

    q := question.QuestionObject{}
    for {
        row := db.QueryRow(query)
        if err = row.Scan(&q.Id, &q.Question); err != nil {
            if  err.Error() == "sql: no rows in result set" {
                continue
            }
            return question.QuestionObject{}, err
        }
        break
    }

    rows, err := db.Query("select * from answers where question_id = $1", q.Id)
    if err != nil {
        return question.QuestionObject{}, err
    }
    defer rows.Close()
    answers := []question.AnswerObject{}
    for rows.Next() {
        var a question.AnswerObject
        err := rows.Scan(&a.QuestionId, &a.Answer, &a.IsCorrect)
        if err != nil {
            return question.QuestionObject{}, err
        }
        answers = append(answers, a)
    }

    q.Answer = answers

    return q, nil
}

func (s *Storage) GetUser(id string) (user.User, error) {
    u := user.User{}
    db, err := s.client.getConnection()
    defer db.Close()
	if err != nil {
        return u, err
	}
    row := db.QueryRow("select * from users where id = $1", id)
    if err = row.Scan(&u.Id, &u.Username); err != nil && err.Error() != "sql: no rows in result set" {
        return u, err
    }
    return u, nil
}

func (s *Storage) SaveGame(q question.QuestionObject, u user.User) error {
    db, err := s.client.getConnection()
	if err != nil {
        return err
	}

    var rightAnswer int
    for i, a := range q.Answer {
        if a.IsCorrect {
            rightAnswer = i + 1
            break
        }
    }

    row := db.QueryRow("select id from game where user_id = $1 and answered = false", u.Id)
    var unfinishedGameId int
    if err = row.Scan(&unfinishedGameId); err != nil && err.Error() != "sql: no rows in result set" {
        return err
    }

    if unfinishedGameId > 0 {
        _, err = db.Exec("update game set rightAnswer = $1 where id = $2", rightAnswer, unfinishedGameId);
    } else {
        _, err = db.Exec("insert into game (user_id, rightAnswer) values ($1, $2)", u.Id, rightAnswer);
    }
	if err != nil {
        return err
	}

    return nil
}

func (s *Storage) CheckAnswer(userId int) (int, error) {
    var correctAnswer int
    db, err := s.client.getConnection()
	if err != nil {
        return correctAnswer, err
	}

    row := db.QueryRow("select rightAnswer from game where user_id = $1 and answered = false", userId)
    err = row.Scan(&correctAnswer)
    return correctAnswer, err
}

func (s *Storage) MarkAnswered(userId int) error {
    db, err := s.client.getConnection()
	if err != nil {
        return err
	}

    _, err = db.Exec("update game set answered = true where user_id = $1", userId);
    return err
}
