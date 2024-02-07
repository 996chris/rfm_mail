package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbHost := EnvGet("DB_HOST", "3.36.102.51")
	dbPort := EnvGet("DB_PORT", "3306")
	dbUser := EnvGet("DB_USER", "cb2")
	dbPassword := EnvGet("DB_PASSWORD", "cb2@readonly")
	dbName := EnvGet("DB_NAME", "katsu")
	sqlString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", sqlString)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	go func() {
		if err := db.Ping(); err != nil {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}()
	repo := NewMysqlRepository(db)
	rfmBuilder := NewRfmBuilder(repo)
	r := gin.Default()
	NewRFMHttpDeliver(r, rfmBuilder)
	if err := r.Run(); err != nil {
		log.Panic(err)
	}
}
