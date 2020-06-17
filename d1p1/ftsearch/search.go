package app

import (
	"database/sql"
	"log"
)

type SearchResult struct {
	Title       string
	Description string
	Id          int
}

func Search(t string) []SearchResult {
	db, err := sql.Open("mysql", "root:123456@tcp(172.17.0.2:3306)/d1p1")

	if err != nil {
		panic(err.Error())
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		// return
	}
	defer db.Close()

	stmtSlct, err := db.Prepare("SELECT id,title,description from articles where MATCH(title,description) AGAINST (? IN NATURAL LANGUAGE MODE)")
	if err != nil {
		panic(err.Error())
	}

	defer stmtSlct.Close()

	var title string
	var description string
	var id int

	rows, err := stmtSlct.Query(t)

	resultSet := []SearchResult{}
	for rows.Next() {

		err := rows.Scan(&id, &title, &description)
		if err != nil {
			log.Fatal(err)
		}
		resultSet = append(resultSet, SearchResult{Id: id, Title: title, Description: description})

	}
	return resultSet
}
