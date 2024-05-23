// cmd/main.go
package main

import (
	"log"

	"github.com/Prashansa-K/serviceCatalog/config"
	"github.com/Prashansa-K/serviceCatalog/internal/db"
	"github.com/Prashansa-K/serviceCatalog/internal/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	_, err := db.GetDB()

	if err != nil {
		log.Fatal("Can not connect to DB: ", err)
	}

	serverConfig := config.GetServerConfig()

	app := echo.New()

	// Register the routes with the database connection
	routes.RegisterRoutes(app)

	app.Logger.Fatal(app.Start(serverConfig.Address))
}
