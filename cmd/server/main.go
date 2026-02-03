package main

import (
	"ecom-appz/internal/config"
	"ecom-appz/internal/db"
	"log"
)

func main() {
	cfg, err := config.SetupEnv()
	
	if err != nil{
		log.Fatalf("config file is not loaded properly %v\n", err)
	}

	db.StartServer(cfg)

	log.Printf("Starting server in %s mode on port %s",
		cfg.AppEnv,
		cfg.AppPort,
	)
}