package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/nullsploit01/go-microservices/authentication/data"
)

const webPort = "80"

var dbConnectCounts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	conn := connectDB()

	if conn == nil {
		log.Panic("Cant connect to DB")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Panicln("Connecting to DB failed")
			dbConnectCounts += 1
		} else {
			log.Println("Connected to DB")
			return conn
		}

		if dbConnectCounts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Connecting to Database....")
		time.Sleep(2 * time.Second)
		continue
	}
}
