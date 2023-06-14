package main

import (
	"database/sql"
	"github/npbbright/futureskill/controller"
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
	createLog := `CREATE TABLE IF NOT EXISTS change_log(
		logid INT AUTO_INCREMENT,
		log_time DATETIME NOT NULL,
		command VARCHAR(45) NOT NULL,
		JSON VARCHAR(256) NOT NULL,
		PRIMARY KEY (logid));`

	if _, err := db.Exec(createTb); err != nil {
		log.Fatal("create table error", err)
	}
	if _, err := db.Exec(createLog); err != nil {
		log.Fatal("create table error", err)
	}
	route := gin.Default()
	route.POST("/createmovie", func(c *gin.Context) {
		var creatMovie service.Movie_Service = service.MovieHandle(*c)
		var movieController controller.ControllerMovie = controller.ServiceHandler(&creatMovie)
		movieController.AddMovie(c)
	})
	route.GET("/getmovie/:get_id", func(ctx *gin.Context) {
		var getMovie service.Movie_Service = service.MovieHandle(*ctx)
		var getController controller.ControllerMovie = controller.ServiceHandler(&getMovie)
		getController.GetMovie(ctx)
	})
	route.PUT("/updaterating/:update_id", func(ctx *gin.Context) {
		var updateRate service.Movie_Service = service.MovieHandle(*ctx)
		var updateController controller.ControllerMovie = controller.ServiceHandler(&updateRate)
		updateController.UpdateRate(ctx)
	})
	route.DELETE("/delete/:delete_id", func(ctx *gin.Context) {
		var deleteService service.Movie_Service = service.MovieHandle(*ctx)
		var deleteController controller.ControllerMovie = controller.ServiceHandler(&deleteService)
		deleteController.DelMovie(ctx)
	})
	route.PATCH("/resotre/:restore_id", func(ctx *gin.Context) {
		var restoreService service.Movie_Service = service.MovieHandle(*ctx)
		var resotreContorller controller.ControllerMovie = controller.ServiceHandler(&restoreService)
		resotreContorller.ResMovie(ctx)
	})

	route.Run(":9090")

}
