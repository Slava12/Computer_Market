package handlefunc

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

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

// InitRandomizer инициализирует генератор случайных чисел
func InitRandomizer() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomNumber(min int, max int) int {
	return min + rand.Intn(max-min)
}

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
		session, _ := store.Get(r, "cookie-name")

		email := r.FormValue("login")
		password := r.FormValue("password")
		name := r.FormValue("name")

		_, err := database.GetUserByEmail(email)
		if err == nil {
			message := "Пользователь с указанным адресом почты уже существует!"
			logger.Warn(err, message)
			errorMessage := errortemplate.Error{Message: message, Link: "/create_account"}
			execute(w, "error.html", errorMessage)
			return
		}
		id, err := database.NewUser(0, false, email, password, name, "")
		if err != nil {
			logger.Warn(err, "Не удалось добавить нового пользователя!")
		}
		logger.Info("Добавление пользователя", id, "прошло успешно.")

		code := getRandomNumber(100000, 999999)
		codeID, err := database.NewCode(code, id)
		if err != nil {
			logger.Warn(err, "Не удалось создать код подтверждения!")
		}
		logger.Info("Код подтверждения ", codeID, " успешно создан.")
		body := "Ваш код активации: " + strconv.Itoa(code)
		post.SendMail(email, "Подтверждение регистрации на сайте интернет-магазина", body)
		session.Values["authenticated"] = true
		session.Values["login"] = email
		session.Save(r, w)
		http.Redirect(w, r, "/confirm_account", 302)
	}
}

func confirmAccount(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "cookie-name")
	email, _ := session.Values["login"].(string)
	isLogged, _ := session.Values["authenticated"].(bool)
	if isLogged == false {
		http.Redirect(w, r, "/index", 302)
		return
	}
	user, _ := database.GetUserByEmail(email)
	if user.Confirmed == true {
		http.Redirect(w, r, "/index", 302)
		return
	}

	if r.Method == "GET" {
		execute(w, "confirm_account.html", email)
	}
	if r.Method == "POST" {

		inputCode := r.FormValue("code")

		code, errGet := database.GetCodeByUserID(user.ID)
		if errGet != nil {
			logger.Warn(errGet, "Не удалось получить данные о коде подтверждения!")
		}
		logger.Info("Данные о коде подтверждения получены успешно.")
		if strconv.Itoa(code.Code) != inputCode {
			message := "Введён неверный код!"
			logger.Info(message)
			errorMessage := errortemplate.Error{Message: message, Link: "/confirm_account"}
			execute(w, "error.html", errorMessage)
			return
		}

		err := database.UpdateUserConfirmed(user.ID, true)
		if err != nil {
			logger.Warn(err, "Не удалось обновить запись пользователя ", user.ID, "!")
		} else {
			logger.Info("Запись пользователя ", user.ID, " обновлена успешно.")
		}

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
