package main

import (
	"log"
	"os"
	"strings"
	"time"
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/MartinKacani/hospital-webapi/api"
	"github.com/MartinKacani/hospital-webapi/internal/hospital_wl"
	"github.com/MartinKacani/hospital-webapi/internal/db_service"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("HOSPITAL_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("HOSPITAL_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") {
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// Setup database service for Department documents
	dbService := db_service.NewMongoService[hospital_wl.Department](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	hospital_wl.SeedDepartments(context.Background(), dbService)
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	// Request routings
	handleFunctions := &hospital_wl.ApiHandleFunctions{
		DepartmentsAPI:  hospital_wl.NewDepartmentsApi(),
		ReservationsAPI: hospital_wl.NewReservationsApi(),
		StaysAPI:        hospital_wl.NewStaysApi(),
	}
	hospital_wl.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
