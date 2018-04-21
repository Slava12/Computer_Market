package handlefunc

import (
	"net/http"
	"strconv"

	"github.com/Slava12/Computer_Market/database"
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
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "category.html", category)
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
		err := database.UpdateCategory(result.ID, result.ParentID, result.Name, result.Features)
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
		/*lol, errM := json.Marshal(result.Features) // пример
		if errM != nil {
			fmt.Println("Err Marshal")
		}
		kek := make([]database.Feature, len(result.Features))
		errU := json.Unmarshal(lol, &kek)
		if errU != nil {
			fmt.Println("Err Unmarshal")
		}*/
		id, errAdd := database.NewCategory(result.ParentID, result.Name, result.Features)
		if errAdd != nil {
			logger.Warn(errAdd, "Не удалось добавить новую категорию!")
		} else {
			logger.Info("Добавление категории ", id, " прошло успешно.")
		}
		http.Redirect(w, r, "/edit/categories", 302)
	}
}
