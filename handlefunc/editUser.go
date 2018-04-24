package handlefunc

import (
	"net/http"
	"strconv"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

func users(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUsers()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список пользователей!")
	} else {
		logger.Info("Список пользователей получен успешно.")
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "users.html", users)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func showUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstring := vars["id"]
	userID, errString := strconv.Atoi(userIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
	}
	user, err := database.GetUser(userID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о пользователе ", userID, "!")
	} else {
		logger.Info("Данные о пользователе ", userID, " получены успешно.")
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "user.html", user)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		result := database.User{}
		result.ID, _ = strconv.Atoi(r.FormValue("id"))
		result.AccessLevel, _ = strconv.Atoi(r.FormValue("accessLevel"))
		result.Login = r.FormValue("login")
		result.Password = r.FormValue("password")
		result.Email = r.FormValue("email")
		result.FirstName = r.FormValue("firstName")
		result.SecondName = r.FormValue("secondName")
		err := database.UpdateUser(result.ID, result.AccessLevel, result.Login, result.Password, result.Email, result.FirstName, result.SecondName)
		if err != nil {
			logger.Warn(err, "Не удалось обновить запись пользователя ", result.ID, "!")
		} else {
			logger.Info("Запись пользователя ", result.ID, " обновлена успешно.")
		}
		http.Redirect(w, r, "/edit/users", 302)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "add_user.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		result := database.User{}
		result.AccessLevel, _ = strconv.Atoi(r.FormValue("accessLevel"))
		result.Login = r.FormValue("login")
		result.Password = r.FormValue("password")
		result.Email = r.FormValue("email")
		result.FirstName = r.FormValue("firstName")
		result.SecondName = r.FormValue("secondName")
		id, err := database.NewUser(result.AccessLevel, result.Login, result.Password, result.Email, result.FirstName, result.SecondName)
		if err != nil {
			logger.Warn(err, "Не удалось добавить нового пользователя!")
		} else {
			logger.Info("Добавление пользователя", id, "прошло успешно.")
		}
		http.Redirect(w, r, "/edit/users", 302)
	}
}