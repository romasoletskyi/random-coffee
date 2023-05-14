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

type Database struct {
	db *sql.DB
}

func CreateDatabase(ctx context.Context) (Database, error) {
	db, err := sql.Open("pgx", "user=postgres password=admin host=localhost port=5432 database=postgres sslmode=disable")
	if err != nil {
		return Database{}, err
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS forms(id serial primary key, ts timestamptz, name text, email text,
																   contact text, bio text, target text,
																   latitude float, longitude float, radius float,
																   time text, language text);`)
	if err != nil {
		return Database{db}, err
	}

	return Database{db}, nil
}

func (d *Database) AddUserForm(ctx context.Context, form UserForm) error {
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

func (d *Database) Close() error {
	return d.db.Close()
}
