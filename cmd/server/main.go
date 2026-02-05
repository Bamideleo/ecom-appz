package main

import (
	"ecom-appz/internal/config"
	"ecom-appz/internal/db"
	_"ecom-appz/internal/logger"
	_"ecom-appz/internal/middleware"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.SetupEnv()
	// logger := logger.New()
	
	if err != nil{
		log.Fatalf("config file is not loaded properly %v\n", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("OK"))
	})

	// Middleware chain
	// handler := middleware.Recovery(logger)(
	// 	middleware.RequestLogger(logger)(
	// 		mux,
	// 	),
	// )

	// server := &http.Server{
	// 	Handler: handler,
	// }

	db.StartServer(cfg)

	log.Printf("Starting server in %s mode on port %s",
		cfg.AppEnv,
		cfg.AppPort,
	)
}