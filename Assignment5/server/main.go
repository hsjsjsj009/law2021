package main

import (
	"Assignment5/server/model"
	readServ "Assignment5/server/read"
	"Assignment5/server/write"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)



func main() {

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5433"
	}

	dsn := fmt.Sprintf("host=%s user=user password=password dbname=backend port=%s sslmode=disable TimeZone=Asia/Jakarta", host, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(20)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	err = db.AutoMigrate(&model.Temp{})
	if err != nil {
		log.Fatal(err.Error())
	}

	writeServer := write.Init(db)
	readServer := readServ.Init(db)

	waitChan := make(chan bool)

	go func() {
		if err := writeServer.Listen(":8080"); err != nil {
			log.Fatal(err.Error())
		}
	}()

	go func() {
		if err := readServer.Listen(":8081"); err != nil {
			log.Fatal(err.Error())
		}
	}()

	<-waitChan
}
