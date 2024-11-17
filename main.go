package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tanaphonble/golang-htmx/app/dog"
)

func main() {
	dbConnection, err := sql.Open("sqlite3", "./dog.db")
	if err != nil {
		log.Fatalf("failed to open sqlite db, error: %s", err)
	}
	defer dbConnection.Close()

	dogDatabase := dog.NewDogDatabase(dbConnection)

	r := gin.Default()
	dog.RegisterHandler(r, dogDatabase)

	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")

	r.GET("/version", func(c *gin.Context) {
		c.String(200, "1.0.0")
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "dog-crud.html", gin.H{
			"title": "Golang X HTMX",
		})
	})

	// r.GET("/table-rows", func(ctx *gin.Context) {

	// })

	r.Run() // listen and serve on 0.0.0.0:8080
}
