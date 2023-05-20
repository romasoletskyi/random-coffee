package app

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	"github.com/romasoletskyi/random-coffee/internal/user"
	"github.com/sirupsen/logrus"
)

func Initialize(log string, databaseConstructor func(context.Context) (*sql.DB, error)) (*os.File, *sql.DB) {
	flag.StringVar(&user.Username, "user", "", "smtp username")
	flag.StringVar(&user.Password, "pass", "", "smtp password")
	flag.Parse()

	file, err := os.OpenFile(log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(file)
	logrus.SetLevel(logrus.TraceLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := databaseConstructor(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	return file, db
}
