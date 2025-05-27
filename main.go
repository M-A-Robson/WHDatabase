package main

import (
	"fmt"
	"internal/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitialiseDatabase()
	router := gin.Default()
	router.GET("/models", getAllModels)
	router.GET("/models/id", getModel)
	router.POST("/models", postModels)
	router.PUT("/models", editModel)
	router.Run("localhost:8000")
}

func getAllModels(context *gin.Context) {
	models, err := db.GetAllModels()
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			"error retrieving models")
		return
	}
	context.IndentedJSON(http.StatusOK, models)
}

func getModel(context *gin.Context) {
	id := context.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			fmt.Sprintf("could not convert id: %s to int", id))
		return
	}

	model, err := db.GetModel(i)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to retrieve unit with id %s", id))
		return
	}
	context.IndentedJSON(http.StatusOK, model)
}

func postModels(context *gin.Context) {
	var newUnit db.Model

	if err := context.BindJSON(&newUnit); err != nil {
		context.JSON(
			http.StatusBadRequest,
			"Failed to parse request body data")
		return
	}

	modelId, err := addModelToDatabase(newUnit)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			err.Error())
		return
	}

	newUnit.ID = modelId

	if err = updateUnitData(newUnit); err != nil {
		context.JSON(
			http.StatusInternalServerError,
			err.Error())
		return
	}

	context.IndentedJSON(http.StatusCreated, newUnit)
}

func editModel(context *gin.Context) {
	var updatedUnit db.Model
	if err := context.BindJSON(&updatedUnit); err != nil {
		context.JSON(
			http.StatusBadRequest,
			"Failed to parse request body data")
		return
	}

	if updatedUnit.ID == 0 {
		modelId, err := addModelToDatabase(updatedUnit)
		if err != nil {
			context.JSON(
				http.StatusInternalServerError,
				err.Error())
			return
		}
		updatedUnit.ID = modelId
	}

	if err := updateUnitData(updatedUnit); err != nil {
		context.JSON(
			http.StatusInternalServerError,
			err.Error())
		return
	}

	context.IndentedJSON(http.StatusAccepted, updatedUnit)

}

func addModelToDatabase(newUnit db.Model) (int, error) {
	date, err := db.ParseDate(newUnit.PurchaseDate)
	if err != nil {
		return -1, fmt.Errorf(
			"purchase date format invalid, must be YYYY-MM-DD")
	}

	modelId, err := db.CreateModelEntry(newUnit.Game,
		newUnit.Faction,
		newUnit.UnitName,
		newUnit.UnitSize,
		date)
	if err != nil {
		return -1, fmt.Errorf(
			"could not create model entry")
	}
	return modelId, nil
}

func updateUnitData(unit db.Model) error {
	old_data, err := db.GetModel(unit.ID)
	if err != nil {
		return fmt.Errorf("unit with id %d not found", unit.ID)
	}

	if unit.Faction != old_data.Faction {
		if err := db.UpdateUnitFaction(unit.ID, unit.Faction); err != nil {
			return fmt.Errorf("could not update unit image")
		}
	}

	if unit.UnitName != old_data.UnitName {
		if err := db.UpdateUnitName(unit.ID, unit.UnitName); err != nil {
			return fmt.Errorf("could not update unit image")
		}
	}

	if unit.UnitSize != old_data.UnitSize {
		if err := db.UpdateUnitSize(unit.ID, unit.UnitSize); err != nil {
			return fmt.Errorf("could not update unit image")
		}
	}

	if unit.Points != old_data.Points {
		if err := db.UpdateModelPoints(unit.ID, unit.Points); err != nil {
			return fmt.Errorf("could not update unit points to: %d", unit.Points)
		}
	}

	if unit.PurchaseDate != old_data.PurchaseDate {
		if purchaseDate, err := db.ParseDate(unit.PurchaseDate); err != nil {
			return fmt.Errorf("build date format invalid, must be YYYY-MM-DD")
		} else {
			if err := db.UpdateModelBuildDate(unit.ID, purchaseDate); err != nil {
				return fmt.Errorf("could not update unit build date")
			}
		}
	}

	if unit.BuildDate != old_data.BuildDate {
		if buildDate, err := db.ParseDate(unit.BuildDate); err != nil {
			return fmt.Errorf("build Date format invalid, must be YYYY-MM-DD")
		} else {
			if err := db.UpdateModelBuildDate(unit.ID, buildDate); err != nil {
				return fmt.Errorf("could not update unit build date")
			}
		}
	}

	if unit.PaintedDate != old_data.PaintedDate {
		if paintDate, err := db.ParseDate(unit.PaintedDate); err != nil {
			return fmt.Errorf("painted date format invalid, must be YYYY-MM-DD")
		} else {
			if err := db.UpdateModelPaintedDate(unit.ID, paintDate); err != nil {
				return fmt.Errorf("could not update unit painted date")
			}
		}
	}

	if unit.Image != nil {
		if err := db.UpdateModelImage(unit.ID, unit.Image); err != nil {
			return fmt.Errorf("could not update unit image")
		}
	}

	return nil
}
