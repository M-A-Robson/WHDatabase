/*
Copyright © 2025 Mark Robson https://github.com/M-A-Robson
*/
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

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
	points INT NOT NULL DEFAULT 0,
	purchaseDate TEXT NOT NULL DEFAULT '',
	buildDate TEXT NOT NULL DEFAULT '',
	paintedDate TEXT NOT NULL DEFAULT '',
	image BLOB);`

	_, err = DB.Exec(createTableSql)
	if err != nil {
		log.Fatalf("Error creating table: %q\n", err)
	}
}

func GetAllModels() ([]Model, error) {
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
	purchaseDate Date) (int, error) {
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
	fmt.Printf("unit_id = %d", unit_id)
	return int(unit_id), err
}

func UpdateUnitGame(id int, game string) error {
	_, err := DB.Exec(
		"UPDATE models SET name = ? WHERE id = ?",
		game,
		id)
	return err
}

func UpdateUnitFaction(id int, faction string) error {
	_, err := DB.Exec(
		"UPDATE models SET faction = ? WHERE id = ?",
		faction,
		id)
	return err
}

func UpdateUnitName(id int, name string) error {
	_, err := DB.Exec(
		"UPDATE models SET unitName = ? WHERE id = ?",
		name,
		id)
	return err
}

func UpdateUnitSize(id int, unitSize int) error {
	_, err := DB.Exec(
		"UPDATE models SET unitSize = ? WHERE id = ?",
		unitSize,
		id)
	return err
}

func UpdateModelPurchaseDate(id int, date Date) error {
	datestr := DateToString(date)
	_, err := DB.Exec(
		"UPDATE models SET purchaseDate = ? WHERE id = ?",
		datestr,
		id)
	return err
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
