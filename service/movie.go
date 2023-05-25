package service

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	AddMovie(gin *gin.Context)
}

type Movie struct {
	ImdbID      string  `json:"imdbID"`
	Title       string  `json:"title"`
	Year        int     `json:"year"`
	Rating      float32 `json:"rating"`
	IsSuperHero bool    `json:"isSuperHero"`
	Status      string  `json :"status"`
}

func AddMovie(gin *gin.Context) {
	var a_movie Movie
	body, err := io.ReadAll(gin.Request.Body)
	if err != nil {
		http.Error(gin.Writer, "Failed to Read Request Body", http.StatusBadGateway)

	}
	err = json.Unmarshal(body, &a_movie)
	if err != nil {
		http.Error(gin.Writer, "Failed to parse JSON data", http.StatusBadGateway)
	}
	db, err := sql.Open("ramsql", "goimdb")
	if err != nil {
		log.Fatal("Error : Server")
	}
	defer db.Close()

	addScript := `INSERT INTO goimdb(imdbID,title,year,rating,isSuperHero,status)
	VALUES (?,?,?,?,?,"A")`

	prepare_add, err := db.Prepare(addScript)
	if err != nil {
		log.Fatal("Error : Prepare Stage")
	}
	exec_add, err := prepare_add.Exec(a_movie.ImdbID, a_movie.Title, a_movie.Year, a_movie.Rating, a_movie.IsSuperHero)
	if err != nil {
		log.Fatal("Error : Exec Failed", err, exec_add)

	}
	gin.JSON(http.StatusCreated, "message : Create Movie Sucessful!!")
}

func GetMoive(gin *gin.Context) {
	imdbID := gin.Param("get_id")
	db, err := sql.Open("ramsql", "goimdb")
	if err != nil {
		log.Fatal("Error : Server")
	}
	defer db.Close()
	getScript := db.QueryRow(`SELECT imdbID, title, year, rating, isSuperHero,status 
	FROM goimdb WHERE imdbID=?`, imdbID)
	var getMoive Movie
	err = getScript.Scan(&getMoive.ImdbID, &getMoive.Title, &getMoive.Year, &getMoive.Rating, &getMoive.IsSuperHero, &getMoive.Status)
	switch err {
	case nil:
		if getMoive.Status == "A" {
			gin.JSON(http.StatusOK, getMoive)
			return
		} else {
			gin.JSON(http.StatusNotFound, map[string]string{"message ": "Movie's not found"})
			return
		}
	case sql.ErrNoRows:
		gin.JSON(http.StatusNotFound, map[string]string{"message!": "not found"})
		return
	default:
		gin.JSON(http.StatusInternalServerError, err.Error())
		return
	}

}

func UpdateRating(gin *gin.Context) {
	imdbID := gin.Param("update_id")
	var UpdateRating Movie

	body, err := io.ReadAll(gin.Request.Body)
	if err != nil {
		gin.JSON(http.StatusInternalServerError, err)
		return
	}

	err = json.Unmarshal(body, &UpdateRating)
	if err != nil {
		gin.JSON(http.StatusInternalServerError, err)
		return
	}

	db, err := sql.Open("ramsql", "goimdb")
	if err != nil {
		log.Fatal("Error : Server")
	}
	defer db.Close()
	updateScript := `UPDATE goimdb SET rating = ? WHERE imdbID = ? ;`
	prepareUpdate, err := db.Prepare(updateScript)
	if err != nil {
		log.Fatal("Error : Prepare Stage")
	}
	_, err = prepareUpdate.Exec(UpdateRating.Rating, imdbID)
	switch err {
	case nil:
		gin.JSON(http.StatusOK, "Update Rating Succesful!!")
	default:
		gin.JSON(http.StatusInternalServerError, err)
	}
}
func DelMoive(gin *gin.Context) {
	delImdb := gin.Param("delete_id")
	db, err := sql.Open("ramsql", "goimdb")
	if err != nil {
		log.Fatal("Error : Server")
	}
	defer db.Close()
	delScript := `UPDATE goimdb SET status = ? WHERE imdbID = ?;`
	delPre, err := db.Prepare(delScript)
	if err != nil {
		log.Fatal("Error : Prepare")
	}
	_, err = delPre.Exec("D", delImdb)
	if err != nil {
		log.Fatal("Error : Exec")
	}
	gin.JSON(http.StatusOK, "Delete Sucessful")
}
