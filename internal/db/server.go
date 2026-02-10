package db

import (
	"database/sql"
	"ecom-appz/internal/config"
	"log"
	"time"

	_"github.com/jackc/pgx/v5/stdlib"
)

func StartServer(cfg config.AppConfig) {

	db, err := sql.Open("pgx", cfg.DSN)

	if err !=nil {
		log.Fatal("failed to open database:",err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)


	log.Println("Database connected")



	

	

	if err := db.Ping(); err != nil{
		log.Fatal("failed to connect to database:", err)
	}

}



