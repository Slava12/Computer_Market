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
	http.HandleFunc("/edit/users", users)
	http.HandleFunc("/update_user", updateUser)
	http.HandleFunc("/add_user", addUser)

	http.HandleFunc("/edit/features", features)
	http.HandleFunc("/update_feature", updateFeature)
	http.HandleFunc("/add_feature", addFeature)

	usersPaths := createUsersPaths()
	for i := 0; i < len(usersPaths); i++ {
		http.HandleFunc(usersPaths[i], showUser)
	}

	featuresPaths := createFeaturesPaths()
	for i := 0; i < len(featuresPaths); i++ {
		http.HandleFunc(featuresPaths[i], showFeature)
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

func users(w http.ResponseWriter, r *http.Request) {
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
		err := tpl.ExecuteTemplate(w, "users.html", data)
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

//---------------------------------Features--------------------------------------//
func features(w http.ResponseWriter, r *http.Request) {
	result, err := database.GetAllFeatures()
	if err != nil {
		logger.Error(err, "Не удалось загрузить список характеристик!")
	} else {
		logger.Info("Список характеристик получен успешно.")
	}
	type Data struct {
		ID   int
		Name string
		Link string
	}
	data := make([]Data, len(result))
	for i := 0; i < len(result); i++ {
		data[i].ID = result[i].ID
		data[i].Name = result[i].Name
		data[i].Link = "/edit/features/" + strconv.Itoa(result[i].ID)
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "features.html", data)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func showFeature(w http.ResponseWriter, r *http.Request) {
	splitURL := strings.Split(r.URL.String(), "/")
	featureID, errString := strconv.Atoi(splitURL[3])
	if errString != nil {
		logger.Error(errString, "Не удалось конвертировать строку в число!")
	}
	result, err := database.GetFeature(featureID)
	if err != nil {
		logger.Error(err, "Не удалось получить запись о характеристике!")
	} else {
		logger.Info("Данные о характеристике получены успешно.")
	}
	data := struct {
		Title string
		ID    int
		Name  string
	}{
		Title: string(result.ID),
		ID:    result.ID,
		Name:  result.Name,
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "feature.html", data)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func updateFeature(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "add_feature.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		result := database.Feature{}
		result.ID, _ = strconv.Atoi(r.FormValue("id"))
		result.Name = r.FormValue("name")
		err := database.UpdateFeature(result.ID, result.Name)
		if err != nil {
			logger.Error(err, "Не удалось обновить зхарактеристику!")
		} else {
			logger.Info("Характеристика " + r.FormValue("id") + " обновлена успешно.")
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
		err := database.NewFeature(result.Name)
		if err != nil {
			logger.Error(err, "Не удалось добавить новую!")
		} else {
			logger.Info("Добавление характеристики прошло успешно.")
		}
		http.Redirect(w, r, "/edit/features", 302)
	}
}

func createFeaturesPaths() []string {
	result, err := database.GetAllFeatures()
	if err != nil {
		logger.Error(err, "Не удалось загрузить список характеристик!")
	} else {
		logger.Info("Список характеристик получен успешно.")
	}
	featuresPaths := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		featuresPaths[i] = "/edit/features/" + strconv.Itoa(result[i].ID)
	}
	return featuresPaths
}