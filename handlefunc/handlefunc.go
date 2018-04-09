package handlefunc

import (
	"net/http"
	"text/template"

	"github.com/Slava12/Computer_Market/config"
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

	port := configFile.HTTP.Port

	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	title := "Главная страница"
	isLogged := false
	name := "Святослав"
	data := struct {
		Title    string
		IsLogged bool
		Name     string
	}{
		Title:    title,
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
		index(w, r)
		err := tpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		index(w, r)
		err := tpl.ExecuteTemplate(w, "create_account.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
