package repository

//this postgres.go targets to open a POSTGRE SQL Data Base
//it is the lowest level of app

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	userTable       = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"

	defaultConnTimeout = time.Second
)

var connAttempts = 10

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cnf Config) (*sqlx.DB, error) {
	// db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
	// 	cnf.Host, cnf.Port, cnf.Username, cnf.DBName, cnf.Password, cnf.SSLMode))
	// if err != nil {
	// 	return nil, err
	// }
	for connAttempts > 0 {
		db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cnf.Host, cnf.Port, cnf.Username, cnf.DBName, cnf.Password, cnf.SSLMode))

		if err == nil {
			return db, nil
		}

		logrus.Printf("Postgres is trying to connect, attempts left: %d", connAttempts)
		time.Sleep(defaultConnTimeout)
		connAttempts--
	}

	return nil, errors.New("error connecting database")

	// err = db.Ping()
	// if err != nil {
	// 	return nil, err
	// }

}
