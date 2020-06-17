package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"

	search "./ftsearch"
)

type Login struct {
	User     string `form:"user" json:"user" xml:"user" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type SearchText struct {
	Text string `form:"search" json:"search" xml:"search" binding:"required"`
}

type UserDetail struct {
	FullName string
	Email    string
}

type SearchResult struct {
	Title       string
	Description string
	Id          int
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("public/*")
	r.MaxMultipartMemory = 8 << 20
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"name": "Adesh"})
	})

	r.POST("/", func(c *gin.Context) {

		var json SearchText
		if err := c.ShouldBind(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resultSet := search.Search(json.Text)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"results": resultSet})
	})

	r.POST("/login", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if json.User != "adesh" || json.Password != "123456" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "You are logged in"})
		return
	})

	r.GET("/users", func(c *gin.Context) {
		db, err := sql.Open("mysql", "root:123456@tcp(172.17.0.2:3306)/d1p1")

		if err != nil {
			panic(err.Error())
			// c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			// return
		}
		defer db.Close()

		stmtSlct, err := db.Prepare("SELECT full_name,email from app_users where uid=?")
		if err != nil {
			panic(err.Error())
		}

		defer stmtSlct.Close()

		var email string
		var fullName string
		err = stmtSlct.QueryRow("wwdfwewweerefefeg").Scan(&fullName, &email)

		c.JSON(http.StatusOK, gin.H{"email": email, "full_name": fullName})
	})

	r.GET("/users/:name/:surname/:age", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":     c.Param("name"),
			"surname":  c.Param("surname"),
			"age":      c.Param("age"),
			"location": c.Query("location")})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
