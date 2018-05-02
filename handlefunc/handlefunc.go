package handlefunc

import (
	"net/http"
	"text/template"

	"github.com/Slava12/Computer_Market/config"
	"github.com/Slava12/Computer_Market/database"
	"github.com/gorilla/mux"
)

var (
	tpl         *template.Template
	filesFolder string
)

// InitHTTP инициализирует сетевые функции приложения
func InitHTTP(configFile config.Config) {
	tpl = template.Must(template.ParseGlob("templates/*.html"))

	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	r.PathPrefix("/pictures/").Handler(http.StripPrefix("/pictures/", http.FileServer(http.Dir("./pictures/"))))

	r.HandleFunc("/index", index)
	r.HandleFunc("/login", login)
	r.HandleFunc("/create_account", createAccount)
	r.HandleFunc("/profile", profile)

	r.HandleFunc("/edit", edit)

	r.HandleFunc("/edit/users", users)
	r.HandleFunc("/edit/users/{id}", showUser)
	r.HandleFunc("/update_user", updateUser)
	r.HandleFunc("/add_user", addUser)
	r.HandleFunc("/delete_user", delUser)
	r.HandleFunc("/delete_all_users", delAllUsers)

	r.HandleFunc("/edit/features", features)
	r.HandleFunc("/edit/features/{id}", showFeature)
	r.HandleFunc("/update_feature", updateFeature)
	r.HandleFunc("/add_feature", addFeature)
	r.HandleFunc("/delete_feature", delFeature)
	r.HandleFunc("/delete_all_features", delAllFeatures)

	r.HandleFunc("/edit/categories", categories)
	r.HandleFunc("/edit/categories/{id}", showCategory)
	r.HandleFunc("/update_category", updateCategory)
	r.HandleFunc("/add_category", addCategory)
	r.HandleFunc("/delete_category", delCategory)
	r.HandleFunc("/delete_all_categories", delAllCategories)

	r.HandleFunc("/edit/units", units)
	r.HandleFunc("/edit/units/{id}", showUnit)
	r.HandleFunc("/update_unit", updateUnit)
	r.HandleFunc("/add_unit", addUnit)
	r.HandleFunc("/delete_unit", delUnit)
	r.HandleFunc("/delete_all_units", delAllUnits)

	port := configFile.HTTP.Port

	filesFolder = configFile.Files.Folder

	http.ListenAndServe(":"+port, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		units, _ := database.GetAllUnits()
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "index.html", units)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
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

//-------------------------------------edit------------------------------------------//
func edit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "edit.html", nil)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
