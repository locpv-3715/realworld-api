package routes

import (
	"realworld-api/controllers"
	"realworld-api/middlewares"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"

	_ "realworld-api/docs"

	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("", controllers.Register)
			users.POST("/login", controllers.Login)
		}

		userAuth := api.Group("/user")
		userAuth.Use(middlewares.AuthMiddleware())
		{
			userAuth.GET("", controllers.GetCurrentUser)
			userAuth.PUT("", controllers.UpdateUser)
		}

		profiles := api.Group("/profiles")
		{
			profiles.GET("/:username", middlewares.OptionalAuthMiddleware(), controllers.GetProfile)

			profiles.POST("/:username/follow", middlewares.AuthMiddleware(), controllers.FollowUser)
			profiles.DELETE("/:username/follow", middlewares.AuthMiddleware(), controllers.UnfollowUser)
		}

		api.GET("/tags", controllers.GetTags)

		articles := api.Group("/articles")
		{
			articles.GET("", middlewares.OptionalAuthMiddleware(), controllers.GetArticles)
			articles.GET("/feed", middlewares.AuthMiddleware(), controllers.GetFeed)
			articles.GET("/:slug", middlewares.OptionalAuthMiddleware(), controllers.GetArticle)
			articles.POST("", middlewares.AuthMiddleware(), controllers.CreateArticle)
			articles.PUT("/:slug", middlewares.AuthMiddleware(), controllers.UpdateArticle)
			articles.DELETE("/:slug", middlewares.AuthMiddleware(), controllers.DeleteArticle)

			articles.POST("/:slug/comments", middlewares.AuthMiddleware(), controllers.AddComment)
			articles.DELETE("/:slug/comments/:id", middlewares.AuthMiddleware(), controllers.DeleteComment)
			articles.GET("/:slug/comments", middlewares.OptionalAuthMiddleware(), controllers.GetComments)

			articles.POST("/:slug/favorite", middlewares.AuthMiddleware(), controllers.FavoriteArticle)
			articles.DELETE("/:slug/favorite", middlewares.AuthMiddleware(), controllers.UnfavoriteArticle)
		}
	}
}
