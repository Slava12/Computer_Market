package handlefunc

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/Slava12/Computer_Market/config"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
)

var (
	tpl *template.Template
)

// InitHTTP инициализирует сетевые функции приложения
func InitHTTP(configFile config.Config) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	http.HandleFunc("/index", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/create_account", createAccount)
	http.HandleFunc("/profile", profile)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/edit/users", editUsers)
	http.HandleFunc("/update_user", updateUser)
	http.HandleFunc("/add_user", addUser)

	usersPaths := createUsersPaths()
	for i := 0; i < len(usersPaths); i++ {
		http.HandleFunc(usersPaths[i], showUser)
	}

	port := configFile.HTTP.Port

	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	menu(w, r)
}

func menu(w http.ResponseWriter, r *http.Request) {
	isLogged := false
	name := "Святослав"
	data := struct {
		IsLogged bool
		Name     string
	}{
		IsLogged: isLogged,
		Name:     name,
	}
	if r.Method == "GET" {
		err := tpl.ExecuteTemplate(w, "menu.html", data)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
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
}

func profile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "profile.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func edit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "edit.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

/*func menuEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := tpl.ExecuteTemplate(w, "menu_edit.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}*/

func editUsers(w http.ResponseWriter, r *http.Request) {
	result, err := database.GetAllUsers()
	if err != nil {
		logger.Error(err, "Не удалось загрузить список пользователей!")
	} else {
		logger.Info("Список пользователей получен успешно.")
	}
	type Data struct {
		ID          int
		AccessLevel int
		Login       string
		Password    string
		Email       string
		FirstName   string
		SecondName  string
		Link        string
	}
	data := make([]Data, len(result))
	for i := 0; i < len(result); i++ {
		data[i].ID = result[i].ID
		data[i].AccessLevel = result[i].AccessLevel
		data[i].Login = result[i].Login
		data[i].Password = result[i].Password
		data[i].Email = result[i].Email
		data[i].FirstName = result[i].FirstName
		data[i].SecondName = result[i].SecondName
		data[i].Link = "/edit/users/" + strconv.Itoa(result[i].ID)
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "records_users.html", data)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func createUsersPaths() []string {
	result, err := database.GetAllUsers()
	if err != nil {
		logger.Error(err, "Не удалось загрузить список пользователей!")
	} else {
		logger.Info("Список пользователей получен успешно.")
	}
	usersPaths := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		usersPaths[i] = "/edit/users/" + strconv.Itoa(result[i].ID)
	}
	return usersPaths
}

func showUser(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.String(), "/")
	userID, errString := strconv.Atoi(splitURL[3])
	if errString != nil {
		logger.Error(errString, "Не удалось конвертировать строку в число!")
	}
	result, err := database.GetUser(userID)
	if err != nil {
		logger.Error(err, "Не удалось получить запись о пользователе!")
	} else {
		logger.Info("Данные о пользователе получены успешно.")
	}
	data := struct {
		Title       string
		ID          int
		AccessLevel int
		Login       string
		Password    string
		Email       string
		FirstName   string
		SecondName  string
		Link        string
	}{
		Title:       string(result.ID),
		ID:          result.ID,
		AccessLevel: result.AccessLevel,
		Login:       result.Login,
		Password:    result.Password,
		Email:       result.Email,
		FirstName:   result.FirstName,
		SecondName:  result.SecondName,
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "user.html", data)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "add_user.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
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
			logger.Error(err, "Не удалось обновить запись пользователя!")
		} else {
			logger.Info("Запись пользователя " + r.FormValue("id") + " обновлена успешно.")
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
		err := database.NewUser(result.AccessLevel, result.Login, result.Password, result.Email, result.FirstName, result.SecondName)
		if err != nil {
			logger.Error(err, "Не удалось добавить нового пользователя!")
		} else {
			logger.Info("Добавление пользователя прошло успешно.")
		}
		http.Redirect(w, r, "/edit/users", 302)
	}
}
