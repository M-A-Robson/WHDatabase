package main

import (
	"fmt"
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
			err.Error())
		return
	}
	context.IndentedJSON(http.StatusOK, models)
}

func postModels(context *gin.Context) {
	var newUnit db.Model

	if err := context.BindJSON(&newUnit); err != nil {
		context.JSON(
			http.StatusBadRequest,
			"Failed to parse request body data")
		return
	}

	date, err := db.ParseDate(newUnit.PurchaseDate)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			"Purchase Date format invalid, must be YYYY-MM-DD")
		return
	}

	modelId, err := db.CreateModelEntry(newUnit.Game,
		newUnit.Faction,
		newUnit.UnitName,
		newUnit.UnitSize,
		date)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			"Could not create model entry")
		return
	}

	if newUnit.Points != 0 {
		if err := db.UpdateModelPoints(modelId, newUnit.Points); err != nil {
			context.JSON(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not update unit points to: %d", newUnit.Points))
			return
		}
	}

	if newUnit.BuildDate != "" {
		if buildDate, err := db.ParseDate(newUnit.BuildDate); err != nil {
			context.JSON(
				http.StatusBadRequest,
				"Build Date format invalid, must be YYYY-MM-DD")
			return
		} else {
			if err := db.UpdateModelBuildDate(modelId, buildDate); err != nil {
				context.JSON(
					http.StatusInternalServerError,
					"Could not update unit build date")
				return
			}
		}
	}

	if newUnit.PaintedDate != "" {
		if paintDate, err := db.ParseDate(newUnit.PaintedDate); err != nil {
			context.JSON(
				http.StatusBadRequest,
				"Painted Date format invalid, must be YYYY-MM-DD")
			return
		} else {
			if err := db.UpdateModelPaintedDate(modelId, paintDate); err != nil {
				context.JSON(
					http.StatusInternalServerError,
					"Could not update unit painted date")
				return
			}
		}
	}

	if newUnit.Image != nil {
		if err := db.UpdateModelImage(modelId, newUnit.Image); err != nil {
			context.JSON(
				http.StatusInternalServerError,
				"Could not update unit image")
			return
		}
	}

	newUnit.ID = modelId
	context.IndentedJSON(http.StatusCreated, newUnit)

}
