package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/errortemplate"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/Slava12/Computer_Market/post"
	"github.com/gorilla/sessions"
)

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		session, _ := store.Get(r, "cookie-name")

		email := r.FormValue("login")
		password := r.FormValue("password")

		user, err := database.GetUserByEmail(email)
		if err != nil {
			logger.Warn(err, "Не существует пользователя с указанным адресом почты!")
			message := "Не существует пользователя с указанным адресом почты!"
			errorMessage := errortemplate.Error{Message: message, Link: "/login"}
			execute(w, "error.html", errorMessage)
			return
		}
		if password != user.Password {
			logger.Warn(err, "Пользователь ", user.Email, ": Пароль указан неверно!")
			message := "Пароль указан неверно!"
			errorMessage := errortemplate.Error{Message: message, Link: "/login"}
			execute(w, "error.html", errorMessage)
			return
		}
		session.Values["authenticated"] = true
		session.Values["login"] = user.Email
		session.Values["name"] = user.FirstName
		session.Save(r, w)
		logger.Info("Пользователь ", user.Email, " успешно авторизовался.")
		http.Redirect(w, r, "/index", 302)
	}
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "create_account.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		post.SendMail("slavanosov@yandex.ru", "Подтверждение регистрации на сайте интернет-магазина", "Перейдите по ссылке, чтобы активировать учётную запись: http://78.106.252.55:8080/index")
		http.Redirect(w, r, "/index", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	//if r.Method == "POST" {
	session, _ := store.Get(r, "cookie-name")
	login, _ := session.Values["login"].(string)

	logger.Info("Пользователь ", login, " вышел из аккаунта.")

	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/index", 302)
	//}
}
