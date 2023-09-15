package config

import (
	"database/sql"
	"fmt"

	"github.com/LloydJanseVanRensburg/go_htmx_quiz_app/internal/main/utils"
)

type Category struct {
	ID          int64
	Name        string
	Description string
}

type Question struct {
	ID         int64
	Question   string
	CategoryID int64
	Options    []Option
}

type Option struct {
	ID         int64
	Text       string
	QuestionID int64
}

var db *sql.DB

func ConnectDB() (*sql.DB) {
	var err error
	
	db, err = sql.Open("sqlite3", "app.db")
	if err != nil {
		panic(err)
	}

	return db
}

func DBTablesSetup() {
	query, err := utils.LoadQuery("createTables.sql")
	if err != nil {
		fmt.Println("Error dbTablesSetup() LoadQuery", err)
		panic(err)
	}

	_, err = db.Exec(query);
	if err != nil {
		fmt.Println("Error dbTablesSetup() exec query", err)
		panic(err)
	}
}

func GetCategories() ([]Category, error) {
	query, err := utils.LoadQuery("getCategories.sql")

	if err != nil {
		fmt.Println("Load query error")
		return nil, err
	}

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var categories []Category

	for rows.Next() {
		var category Category
		rows.Scan(&category.ID, &category.Name, &category.Description)
		
		categories = append(categories, category)
	}
	
	return categories, nil
}

func GetQuestions(categoryId int64) ([]Question, error) {
	questionQuery, err := utils.LoadQuery("getQuestionsByCategoryId.sql")
	if err != nil {
		fmt.Println("Load query error", err)
		return nil, err
	}

	rows, err := db.Query(questionQuery, categoryId)
	if err != nil {
		fmt.Println("DB questions query error", err)
		return nil, err
	}

	var questions []Question

	for rows.Next() {
		var question Question
		rows.Scan(&question.ID, &question.Question, &question.CategoryID)

		options,err := GetOptions(question.ID)
		if err != nil {
			fmt.Println("Error getting question options")
			return nil, err
		}

		question.Options = options

		questions = append(questions, question)
	}

	return questions, nil
}

func GetOptions(questionId int64) ([]Option, error) {
	optionQuery, err := utils.LoadQuery("getOptionsByQuestionId.sql")
	if err != nil {
		fmt.Println("Load query error", err)
		return nil, err
	}

	optionRows, err := db.Query(optionQuery, questionId)
	if err != nil {
		fmt.Println("DB options query error", err)
		return nil, err
	}

	var options []Option

	for optionRows.Next() {
		var option Option
		optionRows.Scan(&option.ID, &option.Text, &option.QuestionID)
		options = append(options, option)
	}

	return options, nil
}

func GetQuestionById(questionId int64) (Question, error) {
	questionQuery, err := utils.LoadQuery("getQuestionById.sql")
	if err != nil {
		fmt.Println("Load query error", err)
		return Question{}, err
	}

	optionRow := db.QueryRow(questionQuery, questionId)

	var question Question

	err = optionRow.Scan(&question.ID, &question.Question, &question.CategoryID)
	if err != nil {
		fmt.Println("DB question query error", err)
		return Question{}, err
	}

	options, err := GetOptions(question.ID)
	if err != nil {
		fmt.Println("DB question options query error", err)
		return Question{}, err
	}

	question.Options = options

	return question, nil
}

func GetOptionById(optionId int64) (Option, error) {
	optionQuery, err := utils.LoadQuery("getOptionById.sql")
	if err != nil {
		fmt.Println("Load query error", err)
		return Option{}, err
	}

	optionRow := db.QueryRow(optionQuery, optionId)

	var option Option

	err = optionRow.Scan(&option.ID, &option.Text, &option.QuestionID)
	if err != nil {
		fmt.Println("DB option query error", err)
		return Option{}, err
	}

	return option, nil
}

func GetQuestionOptionMatch(questionId, optionId int64) (bool, error) {
	answerQuery, err := utils.LoadQuery("getQuestionOptionMatch.sql")
	if err != nil {
		fmt.Println("Load query error", err)
		return false, err
	}

	answerRows, err := db.Query(answerQuery, questionId, optionId)
	if err != nil {
		fmt.Println("DB answer query error", err)
		return false, err
	}

	var answers []int64

	for answerRows.Next() {
		var answer int64
		answerRows.Scan(&answer)
		answers = append(answers, answer)
	}

	if len(answers) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
