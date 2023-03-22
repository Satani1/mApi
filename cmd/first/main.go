package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"mApi/pkg/models/mysql"
	"net/http"
	"os"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	socialDB *mysql.UserModel
}

const addr string = "localhost:9001"
const dsn string = "root:12345678@/socialDB?parseTime=true"

func main() {
	//logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERORR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//database
	socialDB, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer socialDB.Close()
	//first model
	App := &Application{
		errorLog: errorLog,
		infoLog:  infoLog,
		socialDB: &mysql.UserModel{DB: socialDB},
	}

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errorLog,
		Handler:  App.Routes(),
	}

	//launch
	infoLog.Printf("Launching server on %s", addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
