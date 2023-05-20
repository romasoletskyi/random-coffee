package main

import (
	_ "embed"
	"net/http"

	"github.com/romasoletskyi/random-coffee/internal/app"
	"github.com/romasoletskyi/random-coffee/internal/data"
	"github.com/romasoletskyi/random-coffee/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	file, db := app.Initialize("main-log", data.CreateRawDatabase)
	defer file.Close()
	defer func() { _ = db.Close() }()

	http.HandleFunc("/submit", func(w http.ResponseWriter, req *http.Request) {
		server.Submit(data.CreateFormDatabase(db), w, req)
	})
	http.HandleFunc("/feedback", func(w http.ResponseWriter, req *http.Request) {
		server.Feedback(data.CreateFeedbackDatabase(db), w, req)
	})

	err := http.ListenAndServeTLS(":443", "cert/certificate.crt", "cert/private.key", nil)
	logrus.Fatal(err)
}
