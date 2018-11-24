package handlefunc

import (
	"net/http"
	"strconv"

	"strings"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/errortemplate"
	"github.com/Slava12/Computer_Market/files"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

func units(w http.ResponseWriter, r *http.Request) {
	units, err := database.GetAllUnits()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список товаров!")
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "units.html", units)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func editUnit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unitIDstring := vars["id"]
	unitID, errString := strconv.Atoi(unitIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
		return
	}
	unit, err := database.GetUnit(unitID)
	if err != nil {
		logger.Warn(err, "Не удалось получить данные о товаре ", unitID, "!")
		return
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
	unitID, errString := strconv.Atoi(r.FormValue("id"))
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
		return
	}
	unit, err := database.GetUnit(unitID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о товаре ", unitID, "!")
		return
	}
	if r.Method == "POST" {
		unit.Name = r.FormValue("name")
		unit.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
		unit.Price, _ = strconv.Atoi(r.FormValue("price"))
		unit.Discount, _ = strconv.Atoi(r.FormValue("discount"))
		err := database.UpdateUnit(unit.ID, unit.Name, unit.CategoryID, unit.Quantity, unit.Price, unit.Discount, unit.Popularity, unit.Features, unit.Pictures)
		if err != nil {
			logger.Warn(err, "Не удалось обновить товар ", unit.ID, "!")
			return
		}
		logger.Info("Товар ", unit.ID, " обновлён успешно.")
		http.Redirect(w, r, "/edit/units", 302)
	}
}

func addUnit(w http.ResponseWriter, r *http.Request) {
	categories, err := database.GetAllCategories()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список категорий!")
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "add_unit.html", categories)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		result := database.Unit{}
		result.Name = r.FormValue("name")
		result.CategoryID, _ = strconv.Atoi(r.FormValue("category"))
		result.Quantity, _ = strconv.Atoi(r.FormValue("quantity"))
		result.Price, _ = strconv.Atoi(r.FormValue("price"))
		result.Discount, _ = strconv.Atoi(r.FormValue("discount"))
		result.Popularity = 0
		features := r.FormValue("features")
		if features != "" {
			arrayString := strings.Split(features, ";")
			result.Features = make([]database.FeatureUnit, len(arrayString))
			for i := 0; i < len(arrayString); i++ {
				res := strings.Split(arrayString[i], "_")
				result.Features[i].Name = res[0]
				result.Features[i].Value = res[1]
			}
		}
		id, errAdd := database.NewUnit(result.Name, result.CategoryID, result.Quantity, result.Price, result.Discount, result.Popularity, result.Features, result.Pictures)
		if errAdd != nil {
			logger.Warn(errAdd, "Не удалось добавить новый товар!")
			message := errortemplate.GenerateMessage(errAdd)
			errorMessage := errortemplate.Error{Message: message, Link: "/add_unit"}
			err := tpl.ExecuteTemplate(w, "error.html", errorMessage)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		logger.Info("Добавление товара ", id, " прошло успешно.")
		filePath := filesFolder + strconv.Itoa(id) + "/"
		files.CreateDirectory(filePath)
		numberOfPictures, _ := strconv.Atoi(r.FormValue("pictures"))
		if numberOfPictures != 0 {
			result.Pictures = make([]string, numberOfPictures)
			for i := 0; i < numberOfPictures; i++ {
				name := "file" + strconv.Itoa(i)
				_, fileHeader, err := r.FormFile(name)
				if err != nil {
					logger.Warn(err, "Ошибка получения файла!")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				fileName := strconv.Itoa(i)
				fileName = files.Save(filePath, fileHeader, fileName)
				result.Pictures[i] = fileName
			}
		}
		errUpdate := database.UpdateUnit(id, result.Name, result.CategoryID, result.Quantity, result.Price, result.Discount, result.Popularity, result.Features, result.Pictures)
		if errUpdate != nil {
			logger.Warn(errUpdate, "Не удалось обновить информацию о товаре ", id, "!")
			return
		}
		logger.Info("Обновление информации о товаре ", id, " прошло успешно.")
		http.Redirect(w, r, "/edit/units", 302)
	}
}

func delUnit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		unitID, _ := strconv.Atoi(r.FormValue("id"))
		err := database.DelUnit(unitID)
		if err != nil {
			logger.Warn(err, "Не удалось удалить запись о товаре ", unitID, "!")
			return
		}
		logger.Info("Удаление записи о товаре ", unitID, " прошло успешно.")
		path := "pictures/" + strconv.Itoa(unitID) + "/"
		files.RemoveAllFiles(path)
		http.Redirect(w, r, "/edit/units", 302)
	}
}

func delAllUnits(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := database.DelAllUnits()
		if err != nil {
			logger.Warn(err, "Не удалось удалить все записи о товарах!")
			return
		}
		logger.Info("Удаление всех записей о товарах прошло успешно.")
		path := "pictures/"
		files.RemoveAllFiles(path)
		files.CreateDirectory(filesFolder)
		http.Redirect(w, r, "/edit/units", 302)
	}
}
