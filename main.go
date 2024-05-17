package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"os"
	"strings"
	"test_project/rates"
	"test_project/settings"
	"test_project/subscriptions"
)

func main() {
	settingsData, err := os.ReadFile("appsettings.json")
	if err != nil {
		panic(err)
	}

	// null symbols in the beginning of the file...
	startFrom := strings.Index(string(settingsData), "{")
	appSettings := settings.AppSettings{}
	err = json.Unmarshal(settingsData[startFrom:], &appSettings)
	if err != nil {
		panic(err)
	}

	connectionFromEnv := os.Getenv("CONNECTION_STRING")
	if connectionFromEnv != "" {
		appSettings.ConnectionString = connectionFromEnv
	}

	app := fiber.New()
	nbuRatesService := rates.NewService(appSettings)
	nbuRatesHandler := rates.NewHandler(nbuRatesService)

	subscriptionsRepository := subscriptions.NewRepository(appSettings.ConnectionString)
	subscriptionsService := subscriptions.NewService(subscriptionsRepository)
	subscriptionsHandler := subscriptions.NewHandler(subscriptionsService)

	cronService := subscriptions.NewCronService(subscriptionsRepository, appSettings, nbuRatesService)

	app.Get("/api/rate", nbuRatesHandler.GetNbuRate)
	app.Post("/api/subscribe", subscriptionsHandler.AddSubscription)

	cronService.Start()
	defer cronService.Stop()
	err = app.Listen(appSettings.Port)
	if err != nil {
		panic(err)
	}
}
