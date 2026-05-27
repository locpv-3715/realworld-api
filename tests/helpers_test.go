package tests

import (
	"realworld-api/routes"
	"realworld-api/testhelpers"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router)
	return router
}

func setup() func() {
	db := testhelpers.SetupTestDB()
	return func() {
		testhelpers.TeardownTestDB(db)
	}
}
