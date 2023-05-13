package main

import (
	"context"
	_ "embed"
	"flag"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/romasoletskyi/random-coffee/internal/data"
	"github.com/romasoletskyi/random-coffee/internal/server"
	"github.com/romasoletskyi/random-coffee/internal/user"
	"github.com/sirupsen/logrus"
)

func main() {
	flag.StringVar(&user.Username, "user", "", "aws ses username")
	flag.StringVar(&user.Password, "pass", "", "aws ses password")
	flag.Parse()

	file, err := os.OpenFile("log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()
	logrus.SetOutput(file)
	logrus.SetLevel(logrus.TraceLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := data.CreateDatabase(ctx)
	if err != nil {
		logrus.Fatal(err)
	}

	http.HandleFunc("/submit", func(w http.ResponseWriter, req *http.Request) {
		server.Submit(db, w, req)
	})

	err = http.ListenAndServeTLS(":443", "cert/certificate.crt", "cert/private.key", nil)
	logrus.Fatal(err)
}
