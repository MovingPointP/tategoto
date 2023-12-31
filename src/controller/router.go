package controller

import (
	"net/http"
	"tategoto/repository"
	"tategoto/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var serviceInstance service.Services

func GetRouter(db *gorm.DB) *gin.Engine {
	//engine作成
	r := gin.Default()

	//instance作成
	repositoryInstance := repository.New(*db)
	serviceInstance = service.New(repositoryInstance)

	//routing
	auth := r.Group("/api")
	{
		auth.POST("/signup", signup)
		auth.POST("/login", login)
	}

	api := r.Group("/api")
	api.Use(tokenRequired()) //事前・事後処理
	{
		api.GET("/users/:id", getUserByID)
		api.GET("/users", getUsers)
		api.GET("/posts/:id", getPostByID)
		api.GET("/posts", getPosts)

		postHasID := api.Group("/")
		postHasID.Use(compareTokenAndPost())
		{
			postHasID.POST("/posts", createPost)
		}
	}

	//routingなし
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "404:NOT FOUND")
	})

	return r
}
