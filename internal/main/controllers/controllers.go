package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/LloydJanseVanRensburg/go_htmx_quiz_app/internal/main/config"
	"github.com/gin-gonic/gin"
)

type QuestionAnswer struct {
	Question config.Question
	IsCorrect bool
}

func GetIndex(c *gin.Context) {
	categories, err := config.GetCategories()
	if err != nil {
		fmt.Println("Error getting categories")
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"categories": categories,
	})
}

func GetQuestionsPage(c *gin.Context) {
	categoryId, err := strconv.ParseInt(c.Param("categoryId"), 10, 64)
	if err != nil {
		fmt.Println("Category Params parse issue")
		return
	}

	questions, err := config.GetQuestions(categoryId)
	if err != nil {
		fmt.Println("Error getting questions")
		return
	}

	c.HTML(http.StatusOK, "questions-page.html", gin.H{
		"questions": questions,
	})
}

func PostScore(c *gin.Context) {
	c.Request.ParseForm()
	
	var answers []QuestionAnswer

	for key, value := range c.Request.PostForm {
		var x QuestionAnswer

		questionId, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			fmt.Println("Error parsing question id")
			return
		}

		optionId, err := strconv.ParseInt(value[0], 10, 64)
		if err != nil {
			fmt.Println("Error parsing option id")
			return
		}

		question, err := config.GetQuestionById(questionId)
		if err != nil {
			fmt.Println("Error getting question by id")
			return
		}

		answer, err := config.GetQuestionOptionMatch(questionId, optionId)
		if err != nil {
			fmt.Println("Error getting answer")
			return
		}

		x.Question = question
		x.IsCorrect = answer

		answers = append(answers, x)
	}

	sortByQuestionId := func (i, j int) bool {
		return answers[i].Question.ID < answers[j].Question.ID
	}

	sort.Slice(answers, sortByQuestionId)

	c.HTML(http.StatusOK, "questions-page.html", gin.H{
		"questions": answers,
		"done": true,
	})
}