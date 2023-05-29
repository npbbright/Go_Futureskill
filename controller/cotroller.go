package controller

import (
	"github/npbbright/futureskill/service"

	"github.com/gin-gonic/gin"
)

type controllerService struct {
	movieService service.Movie_Service
}

type ControllerMovie interface {
	UpdateRate(gin *gin.Context)
	AddMovie(gin *gin.Context)
	GetMovie(gin *gin.Context)
	DelMovie(gin *gin.Context)
	ResMovie(gin *gin.Context)
}

func ServiceHandler(movieService *service.Movie_Service) ControllerMovie {
	return &controllerService{
		movieService: *movieService,
	}
}

func (c *controllerService) AddMovie(gin *gin.Context) {
	c.movieService.AddMovie(gin)
}
func (c *controllerService) UpdateRate(gin *gin.Context) {
	imdbID := gin.Param("update_id")
	c.movieService.UpdateRating(gin, imdbID)
}
func (c *controllerService) GetMovie(gin *gin.Context) {
	imdbID := gin.Param("get_id")
	c.movieService.GetMovie(gin, imdbID)
}
func (c *controllerService) DelMovie(gin *gin.Context) {
	del_id := gin.Param("delete_id")
	c.movieService.DeleteMovie(gin, del_id)
}
func (c *controllerService) ResMovie(gin *gin.Context) {
	restore_id := gin.Param("restore_id")
	c.movieService.RestoreMovie(gin, restore_id)
}
