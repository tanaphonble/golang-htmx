package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")

	r.GET("/version", func(c *gin.Context) {
		c.String(200, "1.0.0")
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Golang X HTMX",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
