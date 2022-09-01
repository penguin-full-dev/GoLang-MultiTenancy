package main

import (
	_ "./docs" // docs is generated by Swag CLI, you have to import it.
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {

	// load environment variables from file.
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Show Swagger pages
	if os.Getenv("environment") == "development" && os.Getenv("showSwag") == "true" {
		if err := open("http://localhost:8000/swagger/index.html"); err != nil {
			fmt.Println("Something has stopped swagger pages from being loaded into the browser..")
		}
	}

	// Configure port
	port := ":" + os.Getenv("port")

	if port == ":" {
		port = ":8000"
	}

	// Start database services and load master database.
	startDatabaseServices()

	// init router
	router := gin.Default()

	router.Use(CORSMiddleware())

	// Setting up our routes on the router.

	// Users
	setupUsersRoutes(router)

	// Master Users
	setupMasterUsersRoutes(router)

	// Add routing for swag
	if os.Getenv("environment") == "development" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Starting the router instance
	if err := router.Run(port); err != nil {
		fmt.Print(err)

	}
}

// Helper function that allows us to open a browser dependant on your OS
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.Abort()
			return
		}
		c.Next()
	}
}
