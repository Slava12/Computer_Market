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
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "users.html", users)
	}
}

func showUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstring := vars["id"]
	userID, errString := strconv.Atoi(userIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
		return
	}
	user, err := database.GetUser(userID)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о пользователе ", userID, "!")
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "user.html", user)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		result := database.User{}
		result.ID, _ = strconv.Atoi(r.FormValue("id"))
		result.AccessLevel, _ = strconv.Atoi(r.FormValue("accessLevel"))
		if r.FormValue("confirmed") != "" {
			result.Confirmed = true
		} else {
			result.Confirmed = false
		}
		result.Email = r.FormValue("email")
		result.Password = r.FormValue("password")
		result.FirstName = r.FormValue("firstName")
		result.SecondName = r.FormValue("secondName")
		result.Phone = r.FormValue("phone")
		err := database.UpdateUser(result.ID, result.AccessLevel, result.Confirmed, result.Email, result.Password, result.FirstName, result.SecondName, result.Phone)
		if err != nil {
			logger.Warn(err, "Не удалось обновить запись пользователя ", result.ID, "!")
			return
		}
		logger.Info("Запись пользователя ", result.ID, " обновлена успешно.")
		http.Redirect(w, r, "/edit/users", 302)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "add_user.html", nil)
	}
	if r.Method == "POST" {
		result := database.User{}
		result.AccessLevel, _ = strconv.Atoi(r.FormValue("accessLevel"))
		if r.FormValue("confirmed") != "" {
			result.Confirmed = true
		} else {
			result.Confirmed = false
		}
		result.Email = r.FormValue("email")
		result.Password = r.FormValue("password")
		result.FirstName = r.FormValue("firstName")
		result.SecondName = r.FormValue("secondName")
		result.Phone = r.FormValue("phone")
		id, err := database.NewUser(result.AccessLevel, result.Confirmed, result.Email, result.Password, result.FirstName, result.SecondName, result.Phone)
		if err != nil {
			logger.Warn(err, "Не удалось добавить нового пользователя!")
			return
		}
		logger.Info("Добавление пользователя", id, "прошло успешно.")
		http.Redirect(w, r, "/edit/users", 302)
	}
}

func delUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userID, _ := strconv.Atoi(r.FormValue("id"))
		err := database.DelUser(userID)
		if err != nil {
			logger.Warn(err, "Не удалось удалить запись о пользователе ", userID, "!")
			return
		}
		logger.Info("Удаление записи о пользователе ", userID, " прошло успешно.")
		http.Redirect(w, r, "/edit/users", 302)
	}
}

func delAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := database.DelAllUsers()
		if err != nil {
			logger.Warn(err, "Не удалось удалить всех пользователей!")
			return
		}
		logger.Info("Удаление всех пользователей прошло успешно.")
		http.Redirect(w, r, "/edit/users", 302)
	}
}
