package handlefunc

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/Slava12/Computer_Market/config"
	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/errortemplate"
	"github.com/Slava12/Computer_Market/logger"
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

	r.HandleFunc("/search", search)

	r.HandleFunc("/login", login)
	r.HandleFunc("/create_account", createAccount)
	r.HandleFunc("/confirm_account", confirmAccount)
	r.HandleFunc("/logout", logout)

	r.HandleFunc("/profile", profile)
	r.HandleFunc("/profile/change", changeProfile)
	r.HandleFunc("/basket", basket)
	r.HandleFunc("/orders", orders)

	r.HandleFunc("/add_basket/{id}", addBasket)

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

	r.HandleFunc("/categories", showCategories)

	r.HandleFunc("/categories/processors", showProcessors)
	r.HandleFunc("/categories/processors/{id}", showProcessor)

	r.HandleFunc("/categories/motherboards", showMotherboards)
	r.HandleFunc("/categories/motherboards/{id}", showMotherboard)

	r.HandleFunc("/categories/videocards", showVideocards)
	r.HandleFunc("/categories/videocards/{id}", showVideocard)

	r.HandleFunc("/categories/rams", showRams)
	r.HandleFunc("/categories/rams/{id}", showRam)

	r.HandleFunc("/constructor", showConstructor)

	port := configFile.HTTP.Port

	filesFolder = configFile.Files.Folder

	http.ListenAndServe(":"+port, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	processors := makeData("Процессор", "processors", "add_basket", "Добавить в корзину")
	motherboards := makeData("Материнская плата", "motherboards", "add_basket", "Добавить в корзину")
	videocards := makeData("Видеокарта", "videocards", "add_basket", "Добавить в корзину")
	rams := makeData("Оперативная память", "rams", "add_basket", "Добавить в корзину")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_units.html", processors)
		execute(w, "show_units.html", motherboards)
		execute(w, "show_units.html", videocards)
		execute(w, "show_units.html", rams)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	message := "Приносим свои извинения, работа над страницей ещё не завершена."
	errorMessage := errortemplate.Error{Message: message, Link: "/index"}
	execute(w, "error.html", errorMessage)
}

func menu(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	isLogged, _ := session.Values["authenticated"].(bool)
	email, _ := session.Values["login"].(string)
	user, _ := database.GetUserByEmail(email)
	accessLevel := false
	if user.AccessLevel == 10 && isLogged {
		accessLevel = true
	} else {
		page := r.URL.String()
		if strings.Contains(page, "edit") {
			http.Redirect(w, r, "/index", 302)
			return
		}
	}
	data := struct {
		IsLogged    bool
		Name        string
		AccessLevel bool
	}{
		IsLogged:    isLogged,
		Name:        user.FirstName,
		AccessLevel: accessLevel,
	}
	if r.Method == "GET" {
		err := tpl.ExecuteTemplate(w, "menu.html", data)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func showCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_categories.html", nil)
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

//---------------------------------------------------------------------------//

// Data хранит данные об товаре для показа
type Data struct {
	Picture  string
	LinkUnit string
	Name     string
	Price    int
	Link     string
	Text     string
}

func makeData(categoryName string, categoryLink string, actionLink string, actionText string) []Data {
	categories, err := database.GetAllCategories()
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список категорий!")
	} else {
		logger.Info("Список категорий получен успешно.")
	}
	categoryID := 0
	for i := 0; i < len(categories); i++ {
		if categories[i].Name == categoryName {
			categoryID = categories[i].ID
		}
	}
	units, err := database.GetUnitsByCategoryID(categoryID)
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список товаров категории ", categoryName, "!")
	} else {
		logger.Info("Список товаров категории ", categoryName, " получен успешно.")
	}

	data := make([]Data, len(units))
	for i := 0; i < len(units); i++ {
		data[i].Picture = units[i].Pictures[0]
		data[i].LinkUnit = "/categories/" + categoryLink + "/" + strconv.Itoa(units[i].ID)
		data[i].Name = units[i].Name
		data[i].Price = units[i].Price
		data[i].Link = "/" + actionLink + "/" + strconv.Itoa(units[i].ID)
		data[i].Text = actionText
	}
	return data
}

func execute(w http.ResponseWriter, templateName string, data interface{}) {
	err := tpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
