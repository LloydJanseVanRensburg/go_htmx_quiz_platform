package main

import (
	"github.com/LloydJanseVanRensburg/go_htmx_quiz_app/internal/main/config"
	"github.com/LloydJanseVanRensburg/go_htmx_quiz_app/internal/main/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()
	config.DBTablesSetup()

	r := gin.Default()
	routes.Handler(r)

	err := r.Run("127.0.0.1:3000")
	if err != nil {
		panic(err)
	}
}





