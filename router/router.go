package router

import (
	"final-project/controllers"
	"final-project/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
		userRouter.DELETE("/", middlewares.Authentication(), controllers.UserDelete)
		userRouter.PUT("/", middlewares.Authentication(), controllers.UserUpdate)
	}

	photoRouter := r.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.POST("/", controllers.CreatePhoto)
		photoRouter.GET("/", controllers.GetPhoto)
		photoRouter.PUT("/:photoID", middlewares.PhotoAuthorization(), controllers.UpdatePhoto)
		photoRouter.DELETE("/:photoID", middlewares.PhotoAuthorization(), controllers.DeletePhoto)
	}

	commentRouter := r.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.POST("/", controllers.CreateComment)
		commentRouter.GET("/", controllers.GetComment)
		commentRouter.PUT("/:commentID", middlewares.CommentAuthorization(), controllers.UpdateComment)
		commentRouter.DELETE("/:commentID", middlewares.CommentAuthorization(), controllers.DeleteComment)
	}

	socmedRouter := r.Group("/socialmedias")
	{
		socmedRouter.Use(middlewares.Authentication())
		socmedRouter.POST("/", controllers.CreateSocmed)
		socmedRouter.GET("/", controllers.GetSocmed)
		socmedRouter.PUT("/:socialMediaID", middlewares.SocialMediaAuthorization(), controllers.UpdateSocmed)
		socmedRouter.DELETE("/:socialMediaID", middlewares.SocialMediaAuthorization(), controllers.DeleteSocmed)

	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Server is running properly!!")
	})

	return r
}
