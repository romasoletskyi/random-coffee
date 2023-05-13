package server

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/romasoletskyi/random-coffee/internal/data"
	"github.com/romasoletskyi/random-coffee/internal/user"
	"github.com/sirupsen/logrus"
)

func setCORS(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Max-Age", "300")
}

func Submit(db data.Database, w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		setCORS(w, req)
		w.WriteHeader(http.StatusOK)
		return
	}

	if method := req.Method; method != http.MethodPost {
		logrus.Info("Got ", method, " - abort request")
		http.Error(w, "/submit is POST handle", http.StatusBadRequest)
		return
	}

	if content := req.Header.Get("Content-type"); content != "application/json" {
		logrus.Info("Got content-type ", content, " - abort request")
		http.Error(w, "Content-type is not application/json", http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	var form data.UserForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go db.AddUserForm(ctx, form)
	go user.SendConfirmationMail(form)

	setCORS(w, req)
	w.WriteHeader(http.StatusOK)
}
