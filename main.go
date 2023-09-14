package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tdewolff/parse/v2/strconv"
)

type Category struct {
	ID 			int64
	Name 		string
	Description string
}

type Question struct {
	ID 				int64
	Question 		string
	CategoryID 		int64
	Options 		[]Option
}

type Option struct {
	ID 			int64
	Text 		string
	QuestionID 	int64
}

var db *sql.DB

func main() {
	connectDB()
	defer db.Close()
	dbTablesSetup()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/public", "./public")

	r.GET("/", getIndex)
	r.GET("/questions/:categoryId", getQuestionsPage)
	r.POST("/score", postScore)

	r.Run("127.0.0.1:3000")
}

func connectDB() {
	var err error

	db, err = sql.Open("sqlite3", "app.db")
	if err != nil {
		panic(err)
	}
}

func dbTablesSetup() {
	query, err := loadQuery("createTables.sql")
	if err != nil {
		fmt.Println("Error dbTablesSetup() loadQuery", err)
		panic(err)
	}

	_, err = db.Exec(query);
	if err != nil {
		fmt.Println("Error dbTablesSetup() exec query", err)
		panic(err)
	}
}

func loadQuery(filename string) (string, error) {
	query, err := os.ReadFile("./sql/" + filename)
	if err != nil {
		return "", err
	}

	return string(query), nil
}

func getIndex(c *gin.Context) {
	query, err := loadQuery("getCategories.sql");
	if err != nil {
		fmt.Println("Load query error");
		return
	}

	rows, err := db.Query(query);
	if err != nil {
		fmt.Println(err);
		return
	}

	var categories []Category

	for rows.Next() {
		var category Category
		rows.Scan(&category.ID, &category.Name, &category.Description)
		categories = append(categories, category)
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"categories": categories,
	})
}

func getQuestionsPage(c *gin.Context) {
	categoryId := c.Param("categoryId");

	questionQuery,err := loadQuery("getQuestionsByCategoryId.sql")
	if err != nil {
		fmt.Println("Load query error", err);
		return
	}

	var questions []Question

	rows, err := db.Query(questionQuery, categoryId)
	if err != nil {
		fmt.Println("DB questions query error",err);
		return
	}

	for rows.Next() {
		var question Question
		rows.Scan(&question.ID, &question.Question, &question.CategoryID)

		optionQuery,err := loadQuery("getOptionsByQuestionId.sql")
		if err != nil {
			fmt.Println("Load query error", err);
			return
		}

		optionRows, err := db.Query(optionQuery, question.ID)
		if err != nil {
			fmt.Println("DB options query error",err);
			return
		}

		var options []Option

		for optionRows.Next() {
			var option Option
			optionRows.Scan(&option.ID, &option.Text, &option.QuestionID)
			options = append(options, option)
		}

		question.Options = options

		questions = append(questions, question)
	}

	c.HTML(http.StatusOK, "questions-page.html", gin.H{
		"questions": questions,
	})
}

func postScore(c *gin.Context) {
	c.Request.ParseForm()
	type QuestionAnswer struct {
		questionId int64
		answerId int64
	}

	var answers []QuestionAnswer

	for key , value := range c.Request.PostForm {
		var x QuestionAnswer
		questionId, _ := strconv.ParseInt([]byte(key))
		answerId, _ := strconv.ParseInt([]byte(value[0]))

		x.questionId = questionId
		x.answerId = answerId

		answers = append(answers, x)
	}

	fmt.Println(answers)
}