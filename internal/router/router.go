package router

import (
	"PR-appointer/internal/handler"
	"PR-appointer/internal/middleware"
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(ctx context.Context, db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	//corsConfig.AllowAllOrigins = true
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig), middleware.MetricsMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	teamHandler := handler.NewTeamHandler(ctx, db)
	userHandler := handler.NewUserHandler(ctx, db)
	PRHandler := handler.NewPRHandler(ctx, db)

	api := router.Group("/")
	{
		teams := api.Group("/team")
		{
			teams.POST("/add", teamHandler.AddTeam)
			teams.GET("/get", teamHandler.GetTeam)
		}

		users := api.Group("/users")
		{
			users.POST("/setIsActive", userHandler.SetStatus)
			users.GET("/getReview", userHandler.GetUserReviews)
		}

		PRs := api.Group("/pullRequest")
		{
			PRs.POST("/create", PRHandler.CreatePR)
			PRs.POST("/merge", PRHandler.MergePR)
			PRs.POST("/reassign", PRHandler.ReassignReviewer)
		}
	}

	return router
}
