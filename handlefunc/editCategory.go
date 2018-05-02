package handlefunc

import (
	"net/http"
	"strconv"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/errortemplate"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

func categories(w http.ResponseWriter, r *http.Request) {
	categories, err := database.GetAllCategories()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список категорий!")
	} else {
		logger.Info("Список категорий получен успешно.")
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "categories.html", categories)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func showCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryIDstring := vars["id"]
	categoryID, errString := strconv.Atoi(categoryIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
	}
	category, err := database.GetCategory(categoryID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о категории ", categoryID, "!")
	} else {
		logger.Info("Данные о категории ", categoryID, " получены успешно.")
	}
	features, err := database.GetAllFeatures()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список характеристик!")
	} else {
		logger.Info("Список характеристик получен успешно.")
	}
	type FeatureBool struct {
		Feature database.Feature
		Overlap bool
	}
	type Data struct {
		Category database.Category
		Features []FeatureBool
	}
	data := Data{}
	data.Category = category
	data.Features = make([]FeatureBool, len(features))
	for i := 0; i < len(features); i++ {
		data.Features[i].Feature = features[i]
		for j := 0; j < len(category.Features); j++ {
			if features[i].ID == category.Features[j].ID {
				data.Features[i].Overlap = true
			}
		}
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "category.html", data)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		result := database.Category{}
		result.ID, _ = strconv.Atoi(r.FormValue("id"))
		result.ParentID, _ = strconv.Atoi(r.FormValue("parentID"))
		result.Name = r.FormValue("name")
		features, err := database.GetAllFeatures()
		if err != nil {
			logger.Warn(err, "Не удалось загрузить список характеристик!")
		} else {
			logger.Info("Список характеристик получен успешно.")
		}
		j := 0
		tempFeatures := make([]database.Feature, len(features))
		for i := 0; i < len(features); i++ {
			if r.FormValue("feature"+strconv.Itoa(features[i].ID)) != "" {
				tempFeatures[j] = features[i]
				j++
			}
		}
		result.Features = make([]database.Feature, j)
		for i := 0; i < j; i++ {
			result.Features[i] = tempFeatures[i]
		}
		err = database.UpdateCategory(result.ID, result.ParentID, result.Name, result.Features)
		if err != nil {
			logger.Warn(err, "Не удалось обновить категорию ", result.ID, "!")
		} else {
			logger.Info("Категория ", result.ID, " обновлена успешно.")
		}
		http.Redirect(w, r, "/edit/categories", 302)
	}
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	features, err := database.GetAllFeatures()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список характеристик!")
	} else {
		logger.Info("Список характеристик получен успешно.")
	}
	if r.Method == "GET" {
		menu(w, r)
		err = tpl.ExecuteTemplate(w, "add_category.html", features)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		result := database.Category{}
		result.ParentID, _ = strconv.Atoi(r.FormValue("parentID"))
		result.Name = r.FormValue("name")
		j := 0
		tempFeatures := make([]database.Feature, len(features))
		for i := 0; i < len(features); i++ {
			if r.FormValue("feature"+strconv.Itoa(features[i].ID)) != "" {
				tempFeatures[j] = features[i]
				j++
			}
		}
		result.Features = make([]database.Feature, j)
		for i := 0; i < j; i++ {
			result.Features[i] = tempFeatures[i]
		}
		id, errAdd := database.NewCategory(result.ParentID, result.Name, result.Features)
		if errAdd != nil {
			logger.Warn(errAdd, "Не удалось добавить новую категорию!")
			message := errortemplate.GenerateMessage(errAdd)
			errorMessage := errortemplate.Error{Message: message, Link: "/add_category"}
			err := tpl.ExecuteTemplate(w, "error.html", errorMessage)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		logger.Info("Добавление категории ", id, " прошло успешно.")
		http.Redirect(w, r, "/edit/categories", 302)
	}
}

func delCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		categoryID, _ := strconv.Atoi(r.FormValue("id"))
		err := database.DelCategory(categoryID)
		if err != nil {
			logger.Warn(err, "Не удалось удалить запись о категории ", categoryID, "!")
		} else {
			logger.Info("Удаление записи о категории ", categoryID, " прошло успешно.")
		}
		http.Redirect(w, r, "/edit/categories", 302)
	}
}

func delAllCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := database.DelAllCategories()
		if err != nil {
			logger.Warn(err, "Не удалось удалить все записи о категориях!")
		} else {
			logger.Info("Удаление всех записей о категориях прошло успешно.")
		}
		http.Redirect(w, r, "/edit/categories", 302)
	}
}
