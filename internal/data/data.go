package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type mapInfo struct {
	Lat    float32 `json:"lat"`
	Lng    float32 `json:"lng"`
	Radius float32 `json:"radius"`
}

type timeTable = [][]int

type languageTable = []string

func Serialize(x any) (string, error) {
	var builder strings.Builder
	err := json.NewEncoder(&builder).Encode(x)
	if err != nil {
		return "", nil
	}
	return builder.String(), nil
}

type UserForm struct {
	Name    string        `json:"name"`
	Email   string        `json:"email"`
	Contact string        `json:"contact-info"`
	Bio     string        `json:"bio"`
	Target  string        `json:"searching-for"`
	Map     mapInfo       `json:"map"`
	Time    timeTable     `json:"time"`
	Lang    languageTable `json:"lang"`
}

func GetRawDatabase(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("pgx", "user=postgres password=admin host=localhost port=5432 database=postgres sslmode=disable")
	return db, err
}

func CreateRawDatabase(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("pgx", "user=postgres password=admin host=localhost port=5432 database=postgres sslmode=disable")
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS forms(id serial primary key, ts timestamptz, name text, email text,
									contact text, bio text, target text,
									latitude float, longitude float, radius float,
									time text, language text);`)
	if err != nil {
		return db, err
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS pairs(id1 serial, id2 serial);`)
	if err != nil {
		return db, err
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS feedbacks(id serial primary key, meet bool, satisfaction text, add text);`)
	return db, err
}

type FormDatabase struct {
	db *sql.DB
}

func CreateFormDatabase(db *sql.DB) FormDatabase {
	return FormDatabase{db}
}

func (d *FormDatabase) AddUserForm(ctx context.Context, form UserForm) error {
	time, err := Serialize(form.Time)
	if err != nil {
		return err
	}

	language, err := Serialize(form.Lang)
	if err != nil {
		return err
	}

	_, err = d.db.ExecContext(ctx, `INSERT INTO forms (ts, name, email, contact, bio, target, latitude, longitude, radius, time, language) VALUES (current_timestamp, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		form.Name, form.Email, form.Contact, form.Bio, form.Target,
		form.Map.Lat, form.Map.Lng, form.Map.Radius, time, language)

	return err
}

type PairInfo struct {
	Name    string
	Email   string
	Contact string
	Bio     string
}

type PairForm struct {
	ActionLink string
	Left       PairInfo
	Right      PairInfo
}

func ReversePairForm(form PairForm) PairForm {
	return PairForm{Left: form.Right, Right: form.Left}
}

type PairDatabase struct {
	db *sql.DB
}

func CreatePairDatabase(db *sql.DB) PairDatabase {
	return PairDatabase{db}
}

func (d *PairDatabase) GetPairs(ctx context.Context) ([]PairForm, error) {
	rows, err := d.db.QueryContext(ctx,
		`SELECT
			f1.name AS name1,
			f1.email AS email1,
			f1.contact AS contact1,
			f1.bio AS bio1,
			f2.name AS name2,
			f2.email AS email2,
			f2.contact AS contact2,
			f2.bio AS bio2
		FROM
			forms f1
			JOIN pairs p ON f1.id = p.id1
			JOIN forms f2 ON f2.id = p.id2;`)

	forms := make([]PairForm, 0)
	if err != nil {
		return forms, err
	}

	for rows.Next() {
		var form PairForm
		form.ActionLink = "https://romasoletskyi.github.io/random-coffee/pages/feedback.html"

		if err := rows.Scan(&form.Left.Name, &form.Left.Email, &form.Left.Contact, &form.Left.Bio,
			&form.Right.Name, &form.Right.Email, &form.Right.Contact, &form.Right.Bio); err != nil {
			return forms, err
		}
		forms = append(forms, form)
	}

	return forms, rows.Err()
}

type FeedbackForm struct {
	Meet         bool   `json:"meet"`
	Satisfaction string `json:"satisfaction"`
	Add          string `json:"add"`
}

type FeedbackDatabase struct {
	db *sql.DB
}

func CreateFeedbackDatabase(db *sql.DB) FeedbackDatabase {
	return FeedbackDatabase{db}
}

func (d *FeedbackDatabase) AddFeedbackForm(ctx context.Context, form FeedbackForm) error {
	_, err := d.db.ExecContext(ctx, `INSERT INTO feedbacks (meet, satisfaction, add) VALUES ($1, $2, $3);`,
		form.Meet, form.Satisfaction, form.Add)
	return err
}
