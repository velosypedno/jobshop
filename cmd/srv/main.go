package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/velosypedno/jobshop/internal/db"
	"github.com/velosypedno/jobshop/internal/repo"
)

func main() {
	router := gin.Default()
	dbPath := "./db/main.db"
	migrationsPath := "./db/migrations"

	database, err := db.InitDBAndMigrate(dbPath, migrationsPath)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	strategyRepo := repo.NewStrategyRepo(database)

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/strategies", func(ctx *gin.Context) {
		stub := &repo.Strategy{
			Name:   "Stub Strategy",
			Type:   "test",
			Params: `{"key": "value", "status": "stub"}`,
		}

		id, err := strategyRepo.Create(ctx.Request.Context(), stub)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		stub.ID = id
		ctx.JSON(http.StatusCreated, stub)
	})

	router.GET("/strategies", func(ctx *gin.Context) {
		strategies, err := strategyRepo.GetAll(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, strategies)
	})

	router.Run()
}
