package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"ps_ws_ex/planetscale"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ErrorMessage struct {
	Message string `json:"error"`
}

func getReminderService(context *gin.Context) {
	reminder, err := planetscale.GetReminder(1)
	if err != nil {
		context.JSON(http.StatusOK, ErrorMessage{Message: err.Error()})
	} else {
		context.JSON(http.StatusOK, reminder)
	}
}

func RunWebService(port int) error {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	// Serve static files
	router.Static("/", "./static")

	// Initialize the routes
	router.POST("/get-reminder", getReminderService)

	log.Println("Starting on port", port)

	// Start serving the application
	err := router.Run(fmt.Sprintf(":%d", port))

	return err
}

func main() {
	rand.Seed(time.Now().UnixNano())

	port := os.Getenv("PORT")
	portId := 8080
	if port != "" {
		var err error
		portId, err = strconv.Atoi(port)
		if err != nil || portId < 1 {
			panic("Cannot convert " + port + " to an integer (or it's a bad number).")
		}
	}
	err := RunWebService(portId)
	if err != nil {
		panic(err)
	}
}
