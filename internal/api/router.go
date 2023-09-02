package api

import (
	"net/http"
	"time"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/env"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/postgresdb"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/repository"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func newSegmentsRoute(env *env.Env, db *postgresdb.Store, group *gin.RouterGroup) {
	sr := repository.NewSegmentRepository(db)
	sc := segmentController{
		segmentService: services.NewSegmentService(sr, time.Duration(env.App.Timeout)*time.Second),
		env:            env,
	}
	group.POST("/segments", sc.createSegment)
	group.DELETE("/segments", sc.deleteSegment)
}

func newUserSegmentsRoute(env *env.Env, db *postgresdb.Store, group *gin.RouterGroup) {
	usr := repository.NewUserSegmentRepository(db)
	usc := userSegmentController{
		userSegmentService: services.NewUserSegmentService(usr, time.Duration(env.App.Timeout)*time.Second),
		env:                env,
	}
	group.GET("/user_segments", usc.getUserSegments)
	group.POST("/user_segments", usc.updateUserSegments)
}

func newOperationsRoute(env *env.Env, db *postgresdb.Store, group *gin.RouterGroup) {
	or := repository.NewOperationRepository(db)
	oc := operationController{
		operationService: services.NewOperationService(or, time.Duration(env.App.Timeout)*time.Second),
		env:              env,
	}
	group.GET("/operations", oc.getOperations)
}

func SetupRouter(env *env.Env, db *postgresdb.Store, ge *gin.Engine) {
	rootRouter := ge.Group("")
	rootRouter.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	rootRouter.Static("/static", "./web/static")
	rootRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicApiRouter := ge.Group("/api")
	newSegmentsRoute(env, db, publicApiRouter)
	newUserSegmentsRoute(env, db, publicApiRouter)
	newOperationsRoute(env, db, publicApiRouter)
}
