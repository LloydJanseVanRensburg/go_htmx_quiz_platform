package routes

import (
	"github.com/LloydJanseVanRensburg/go_htmx_quiz_app/internal/main/controllers"
	"github.com/gin-gonic/gin"
)

func Handler(r *gin.Engine) {
	r.LoadHTMLGlob("web/template/*")
	r.Static("/static", "web/static")

	r.GET("/", controllers.GetIndex)
	r.GET("/questions/:categoryId", controllers.GetQuestionsPage)
	r.POST("/score", controllers.PostScore)
}