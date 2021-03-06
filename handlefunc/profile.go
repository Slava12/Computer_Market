package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
)

func profile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	isLogged, _ := session.Values["authenticated"].(bool)
	if isLogged == false {
		http.Redirect(w, r, "/index", 302)
		return
	}
	login, _ := session.Values["login"].(string)
	user, err := database.GetUserByEmail(login)
	if err != nil {
		logger.Warn(err, "Не удалось получить информацию о пользователе!")
		http.Redirect(w, r, "/index", 302)
		return
	}
	if user.Confirmed == false {
		http.Redirect(w, r, "/confirm_account", 302)
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "profile.html", user)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func changeProfile(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	email, _ := session.Values["login"].(string)
	user, err := database.GetUserByEmail(email)
	if err != nil {
		logger.Warn(err, "Не удалось получить запись о пользователе ", email, "!")
		return
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "change_profile.html", user)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		firstName := r.FormValue("firstName")
		secondName := r.FormValue("secondName")
		password := r.FormValue("password")
		phone := r.FormValue("phone")
		err := database.UpdateUser(user.ID, user.AccessLevel, user.Confirmed, user.Email, password, firstName, secondName, phone)
		if err != nil {
			logger.Warn(err, "Не удалось обновить запись пользователя ", user.ID, "!")
			return
		}
		logger.Info("Запись пользователя ", user.ID, " обновлена успешно.")
		http.Redirect(w, r, "/profile", 302)
	}
}
