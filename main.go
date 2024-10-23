package main

import (
	"context"
	"net/http"

	"suvva-geo-ride-service/internal/config"
	records_handlers "suvva-geo-ride-service/internal/georecords/handlers"
	geozones_handlers "suvva-geo-ride-service/internal/geozones/handlers"
	"suvva-geo-ride-service/internal/logger"

	"suvva-geo-ride-service/pkg/services"

	"github.com/gorilla/mux"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to load configuration: %v", err)
	}

	client := services.GetClient()
	defer client.Disconnect(context.Background())

	logger.InfoLogger.Println("Successfully connected to MongoDB")

	nc, err := services.GetNatsConnection()
	if err != nil {
		logger.ErrorLogger.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	js := services.GetJetStreamContext(nc)

	services.SetupSubscriptions(js, nc, client)

	router := mux.NewRouter()
	router.HandleFunc("/geozones/{geozoneId}", geozones_handlers.GetGeozoneById(client)).Methods("GET")
	router.HandleFunc("/geozones/{geozoneId}", geozones_handlers.UpdateGeozone(client)).Methods("PUT")
	router.HandleFunc("/geozones/{geozoneId}", geozones_handlers.DeleteGeozone(client)).Methods("DELETE")
	router.HandleFunc("/geozones", geozones_handlers.CreateGeozone(client)).Methods("POST")

	router.HandleFunc("/records", records_handlers.CreateGeoRecord(client)).Methods("POST")
	router.HandleFunc("/update-record-geozones", records_handlers.UpdateRecordGeozones(client)).Methods("POST")

	logger.InfoLogger.Println("Starting server on :8080")
	logger.ErrorLogger.Fatal(http.ListenAndServe(":8080", router))
}
