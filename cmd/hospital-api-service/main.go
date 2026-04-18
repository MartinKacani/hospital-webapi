package main

import (
    "log"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/MartinKacani/hospital-webapi/api"
    "github.com/MartinKacani/hospital-webapi/internal/hospital_wl"
)

func main() {
    log.Printf("Server started")
    port := os.Getenv("HOSPITAL_API_PORT")
    if port == "" {
        port = "8080"
    }
    environment := os.Getenv("HOSPITAL_API_ENVIRONMENT")
    if !strings.EqualFold(environment, "production") { // case insensitive comparison
        gin.SetMode(gin.DebugMode)
    }
    engine := gin.New()
    engine.Use(gin.Recovery())
    // request routings
    handleFunctions := &hospital_wl.ApiHandleFunctions{
        HospitalConditionsAPI:  hospital_wl.NewHospitalConditionsApi(),
        HospitalWaitingListAPI: hospital_wl.NewHospitalWaitingListApi(),
    }
    hospital_wl.NewRouterWithGinEngine(engine, *handleFunctions)
    engine.GET("/openapi", api.HandleOpenApi)
    engine.Run(":" + port)
}