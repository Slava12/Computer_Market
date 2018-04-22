package handlefunc

import (
	"net/http"
	"strconv"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

func units(w http.ResponseWriter, r *http.Request) {
	units, err := database.GetAllUnits()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список товаров!")
	} else {
		logger.Info("Список товаров получен успешно.")
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "units.html", units)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func showUnit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unitIDstring := vars["id"]
	unitID, errString := strconv.Atoi(unitIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
	}
	unit, err := database.GetUnit(unitID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о товаре ", unitID, "!")
	} else {
		logger.Info("Данные о товаре ", unitID, " получены успешно.")
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "unit.html", unit)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func updateUnit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		result := database.Unit{}
		result.Name = r.FormValue("name")
		result.CategoryID, _ = strconv.Atoi(r.FormValue("categoryID"))
		result.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
		result.Price, _ = strconv.Atoi(r.FormValue("price"))
		result.Discount, _ = strconv.Atoi(r.FormValue("discount"))
		err := database.UpdateUnit(result.ID, result.Name, result.CategoryID, result.Quantity, result.Price, result.Discount, result.Features, result.Pictures)
		if err != nil {
			logger.Warn(err, "Не удалось обновить товар ", result.ID, "!")
		} else {
			logger.Info("Товар ", result.ID, " обновлён успешно.")
		}
		http.Redirect(w, r, "/edit/units", 302)
	}
}

func addUnit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "add_unit.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		result := database.Unit{}
		result.Name = r.FormValue("name")
		result.CategoryID, _ = strconv.Atoi(r.FormValue("categoryID"))
		result.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
		result.Price, _ = strconv.Atoi(r.FormValue("price"))
		result.Discount, _ = strconv.Atoi(r.FormValue("discount"))
		id, errAdd := database.NewUnit(result.Name, result.CategoryID, result.Quantity, result.Price, result.Discount, result.Features, result.Pictures)
		if errAdd != nil {
			logger.Warn(errAdd, "Не удалось добавить новый товар!")
		} else {
			logger.Info("Добавление товара ", id, " прошло успешно.")
		}
		http.Redirect(w, r, "/edit/units", 302)
	}
}
