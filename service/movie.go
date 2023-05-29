package service

import (
	"database/sql"
	"encoding/json"
	"github/npbbright/futureskill/core"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	ImdbID      string  `json:"imdbID"`
	Title       string  `json:"title"`
	Year        int     `json:"year"`
	Rating      float32 `json:"rating"`
	IsSuperHero bool    `json:"isSuperHero"`
	Status      string  `json :"status"`
}
type Movie_Service interface {
	AddMovie(gin *gin.Context)
	GetMovie(gin *gin.Context, imdbID string)
	UpdateRating(gin *gin.Context, imdbID string)
	DeleteMovie(gin *gin.Context, imdbID string)
	RestoreMovie(gin *gin.Context, restoreID string)
}

func MovieHandle(gin gin.Context) Movie_Service {
	var movie core.Movie
	if err := gin.ShouldBind(&movie); err != nil {
		gin.AbortWithStatus(http.StatusBadRequest)
	}
	return &Movie{
		ImdbID:      movie.ImdbID,
		Title:       movie.Title,
		Year:        movie.Year,
		Rating:      movie.Rating,
		IsSuperHero: movie.IsSuperHero,
		Status:      movie.Status,
	}
}

func (a_movie *Movie) AddMovie(gin *gin.Context) {
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

func (getMovie *Movie) GetMovie(gin *gin.Context, imdbID string) {
	db, err := sql.Open("ramsql", "goimdb")
	if err != nil {
		log.Fatal("Error : Server")
	}
	defer db.Close()
	getScript := db.QueryRow(`SELECT imdbID, title, year, rating, isSuperHero,status 
	FROM goimdb WHERE imdbID=?`, imdbID)

	err = getScript.Scan(&getMovie.ImdbID, &getMovie.Title, &getMovie.Year, &getMovie.Rating, &getMovie.IsSuperHero, &getMovie.Status)
	switch err {
	case nil:
		if getMovie.Status == "A" {
			gin.JSON(http.StatusOK, getMovie)
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

func (updateRating *Movie) UpdateRating(gin *gin.Context, imdbID string) {

	body, err := io.ReadAll(gin.Request.Body)
	if err != nil {
		gin.JSON(http.StatusInternalServerError, err)
		return
	}

	err = json.Unmarshal(body, &updateRating)
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
	_, err = prepareUpdate.Exec(updateRating.Rating, imdbID)
	switch err {
	case nil:
		gin.JSON(http.StatusOK, "Update Rating Succesful!!")
	default:
		gin.JSON(http.StatusInternalServerError, err)
	}
}
func (_ *Movie) DeleteMovie(gin *gin.Context, delImdb string) {
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
func (_ *Movie) RestoreMovie(gin *gin.Context, restoreID string) {
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
	_, err = delPre.Exec("A", restoreID)
	if err != nil {
		log.Fatal("Error : Exec")
	}
	gin.JSON(http.StatusOK, "Restore Sucessful")
}
