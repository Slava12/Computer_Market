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

	port := configFile.HTTP.Port

	http.ListenAndServe(":"+port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := tpl.ExecuteTemplate(w, "menu.html", "nil")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
