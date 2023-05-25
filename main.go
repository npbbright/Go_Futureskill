package main

import (
	"database/sql"
	"github/npbbright/futureskill/service"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/proullon/ramsql/driver"
)

func main() {
	db, err := sql.Open("ramsql", "goimdb")
	if err != nil {
		log.Fatal("Error : Server")
	}
	createTb := `
	CREATE TABLE IF NOT EXISTS goimdb (
	id INT AUTO_INCREMENT,
	imdbID TEXT NOT NULL UNIQUE,
	title TEXT NOT NULL,
	year INT NOT NULL,
	rating FLOAT NOT NULL,
	isSuperHero BOOLEAN NOT NULL,
	status TEXT NOT NULL,
	PRIMARY KEY (id)
	);
	`

	if _, err := db.Exec(createTb); err != nil {
		log.Fatal("create table error", err)
	}
	route := gin.Default()
	//route.Use(service.Authenticate())
	route.POST("/createmovie", service.AddMovie)
	route.GET("/getmovie/:get_id", service.GetMoive)
	route.PUT("/updaterating/:update_id", service.UpdateRating)
	route.DELETE("/delete/:delete_id", service.DelMoive)

	route.Run(":9090")

}
