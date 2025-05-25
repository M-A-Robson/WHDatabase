package main

import (
	"encoding/json"
	"internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitialiseDatabase()
	router := gin.Default()
	router.GET("/models", getUnits)
	router.POST("/models", postModels)
	router.Run("localhost:8000")
}

func getUnits(context *gin.Context) {
	models, err := db.GetAllModels()
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			"Failed to retrieve data")
		return
	}
	json_models, err := json.Marshal(models)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			"Failed to marshal data")
		return
	}
	context.IndentedJSON(http.StatusOK, json_models)
}

func postModels(context *gin.Context) {
	var newUnit db.Model

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := context.BindJSON(&newUnit); err != nil {
		context.JSON(
			http.StatusBadRequest,
			"Failed to parse request body data")
		return
	}

	date, err := db.ParseDate(newUnit.PurchaseDate)
	if err != nil {
		println(err)
		context.JSON(
			http.StatusBadRequest,
			"Date format invalid, must be YYYY-MM-DD")
		return
	}

	// Add the new album to the slice.
	db.CreateModelEntry(newUnit.Game,
		newUnit.Faction,
		newUnit.UnitName,
		newUnit.UnitSize,
		date)
	context.IndentedJSON(http.StatusCreated, newUnit)
}
