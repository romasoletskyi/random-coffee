package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

//go:embed confirmation-template
var confirmationEmail string

var port, username, password string

func sendMail(to, subject, text string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "random.coffee.manager@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)

	d := gomail.NewDialer("email-smtp.eu-west-3.amazonaws.com", 587, username, password)
	if err := d.DialAndSend(m); err != nil {
		logrus.Error(err)
	}
}

type mapInfo struct {
	Lat    float32 `json:"lat"`
	Lng    float32 `json:"lng"`
	Radius float32 `json:"radius"`
}

type userForm struct {
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Contact string   `json:"contact-info"`
	Bio     string   `json:"bio"`
	Target  string   `json:"searching-for"`
	Map     mapInfo  `json:"map"`
	Time    [][]int  `json:"time"`
	Lang    []string `json:"lang"`
}

func setCORS(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Max-Age", "300")
}

func submit(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodOptions {
		logrus.Info("Got OPTIONS - set CORS")
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

	var form userForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	logrus.Info(form)

	go func() {
		t, err := template.New("make-letter").Parse(confirmationEmail)
		if err != nil {
			logrus.Error(err)
			return
		}

		var builder strings.Builder
		err = t.Execute(&builder, form)
		if err != nil {
			logrus.Error(err)
			return
		}

		logrus.Info(builder.String())
		sendMail(form.Email, "Submit confirmation", builder.String())
	}()

	setCORS(w, req)
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.StringVar(&username, "user", "", "aws ses username")
	flag.StringVar(&password, "pass", "", "aws ses password")
	flag.Parse()

	file, err := os.OpenFile("log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()
	logrus.SetOutput(file)
	logrus.SetLevel(logrus.TraceLevel)

	http.HandleFunc("/submit", submit)

	err = http.ListenAndServeTLS(":3000", "cert/certificate.crt", "cert/private.key", nil)
	//err = http.ListenAndServe(":3000", nil)
	logrus.Fatal(err)
}
