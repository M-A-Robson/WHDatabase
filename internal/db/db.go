/*
Copyright © 2025 Mark Robson https://github.com/M-A-Robson
*/
package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Model struct {
	ID           int
	Game         string
	Faction      string
	UnitName     string
	UnitSize     int
	Points       int
	PurchaseDate string
	BuildDate    string
	PaintedDate  string
	Image        []byte
}

type Date struct {
	Day   int
	Month int
	Year  int
}

func DateToString(date Date) string {
	return fmt.Sprintf("%04d-%02d-%02d", date.Year, date.Month, date.Day)
}

func ParseDate(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return Date{Day: t.Day(), Month: int(t.Month()), Year: t.Year()}, nil
}

func Today() Date {
	now := time.Now()
	return Date{
		Day:   now.Day(),
		Month: int(now.Month()),
		Year:  now.Year(),
	}
}

var DB *sql.DB

func InitialiseDatabase() {
	var err error
	DB, err = sql.Open("sqlite3", "./models.db")
	if err != nil {
		log.Fatalf("Error opening database: %q\n", err)
	}

	createTableSql := `
	CREATE TABLE IF NOT EXISTS models (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	game TEXT NOT NULL,
	faction TEXT NOT NULL,
	unitName TEXT NOT NULL,
	unitSize INT,
	points INT,
	purchaseDate TEXT NOT NULL DEFAULT '',
	buildDate TEXT NOT NULL DEFAULT '',
	paintedDate TEXT NOT NULL DEFAULT '',
	image BLOB);`

	_, err = DB.Exec(createTableSql)
	if err != nil {
		log.Fatalf("Error creating table: %q\n", err)
	}
}

func GetAllTodos() ([]Model, error) {
	models := []Model{}
	rows, err := DB.Query("SELECT * FROM models")
	if err != nil {
		return models, err
	}
	defer rows.Close()
	for rows.Next() {
		var model Model
		if err := rows.Scan(
			&model.ID,
			&model.Game,
			&model.Faction,
			&model.UnitName,
			&model.UnitSize,
			&model.Points,
			&model.PurchaseDate,
			&model.BuildDate,
			&model.PaintedDate,
			&model.Image); err != nil {
			return models, err
		}
		models = append(models, model)
	}
	return models, nil
}

func GetModel(id int) (*Model, error) {
	row := DB.QueryRow(
		"SELECT * FROM models WHERE id = ?",
		id)
	model := &Model{}
	err := row.Scan(
		&model.ID,
		&model.Game,
		&model.Faction,
		&model.UnitName,
		&model.UnitSize,
		&model.Points,
		&model.PurchaseDate,
		&model.BuildDate,
		&model.PaintedDate,
		&model.Image)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func CreateModelEntry(
	game string,
	faction string,
	unitName string,
	unitSize int,
	purchaseDate Date) (int64, error) {
	date := DateToString(purchaseDate)
	result, err := DB.Exec(
		"INSERT INTO models(game, faction, unitName, unitSize, purchaseDate) VALUES(?,?,?,?,?)",
		game,
		faction,
		unitName,
		unitSize,
		date)
	if err != nil {
		return -1, err
	}
	unit_id, err := result.LastInsertId()
	return unit_id, err
}

func UpdateModelPoints(id int, points int) error {
	_, err := DB.Exec(
		"UPDATE models SET points = ? WHERE id = ?",
		points,
		id)
	return err
}

func UpdateModelPaintedDate(id int, date Date) error {
	datestr := DateToString(date)
	_, err := DB.Exec(
		"UPDATE models SET paintedDate = ? WHERE id = ?",
		datestr,
		id)
	return err
}

func UpdateModelBuildDate(id int, date Date) error {
	datestr := DateToString(date)
	_, err := DB.Exec(
		"UPDATE models SET buildDate = ? WHERE id = ?",
		datestr,
		id)
	return err
}

func UpdateModelImage(id int, image []byte) error {
	_, err := DB.Exec(
		"UPDATE models SET image = ? WHERE id = ?",
		image,
		id)
	return err
}

func DeleteModel(id int) error {
	_, err := DB.Exec(
		"DELETE FROM models WHERE id = ?",
		id)
	return err
}

func EditModelId(id int, new_id int) error {
	_, err := DB.Exec(
		"UPDATE models SET id = ? WHERE id = ?",
		new_id,
		id)
	return err
}

func SetAutoIncrementCounter(new_id int) error {
	_, err := DB.Exec(
		"UPDATE SQLITE_SEQUENCE SET seq = ? WHERE name = 'models'",
		new_id)
	return err
}
