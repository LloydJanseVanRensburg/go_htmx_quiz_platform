package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Category struct {
	ID 			int64
	Name 		string
	Description string
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
		fmt.Println("Hello world");
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