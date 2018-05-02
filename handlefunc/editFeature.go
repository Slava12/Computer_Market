package handlefunc

import (
	"net/http"
	"strconv"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/errortemplate"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

func features(w http.ResponseWriter, r *http.Request) {
	features, err := database.GetAllFeatures()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список характеристик!")
	} else {
		logger.Info("Список характеристик получен успешно.")
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "features.html", features)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func showFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	featureIDstring := vars["id"]
	featureID, errString := strconv.Atoi(featureIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
	}
	feature, err := database.GetFeature(featureID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о характеристике ", featureID, "!")
	} else {
		logger.Info("Данные о характеристике ", featureID, " получены успешно.")
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "feature.html", feature)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func updateFeature(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		result := database.Feature{}
		result.ID, _ = strconv.Atoi(r.FormValue("id"))
		result.Name = r.FormValue("name")
		err := database.UpdateFeature(result.ID, result.Name)
		if err != nil {
			logger.Warn(err, "Не удалось обновить характеристику ", result.ID, "!")
		} else {
			logger.Info("Характеристика ", result.ID, " обновлена успешно.")
		}
		http.Redirect(w, r, "/edit/features", 302)
	}
}

func addFeature(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "add_feature.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		result := database.Feature{}
		result.Name = r.FormValue("name")
		id, errAdd := database.NewFeature(result.Name)
		if errAdd != nil {
			logger.Warn(errAdd, "Не удалось добавить новую характеристику!")
			message := errortemplate.GenerateMessage(errAdd)
			errorMessage := errortemplate.Error{Message: message, Link: "/add_feature"}
			err := tpl.ExecuteTemplate(w, "error.html", errorMessage)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		logger.Info("Добавление характеристики ", id, " прошло успешно.")
		http.Redirect(w, r, "/edit/features", 302)
	}
}

func delFeature(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		featureID, _ := strconv.Atoi(r.FormValue("id"))
		err := database.DelFeature(featureID)
		if err != nil {
			logger.Warn(err, "Не удалось удалить запись о характеристике ", featureID, "!")
		} else {
			logger.Info("Удаление записи о характеристике ", featureID, " прошло успешно.")
		}
		http.Redirect(w, r, "/edit/features", 302)
	}
}

func delAllFeatures(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := database.DelAllFeatures()
		if err != nil {
			logger.Warn(err, "Не удалось удалить все записи о характеристиках!")
		} else {
			logger.Info("Удаление всех записей о характеристиках прошло успешно.")
		}
		http.Redirect(w, r, "/edit/features", 302)
	}
}
