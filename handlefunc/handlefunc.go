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
	r.HandleFunc("/basket", ShowBasket)
	r.HandleFunc("/orders", showOrders)

	r.HandleFunc("/add_basket/{id}", AddUnit)
	r.HandleFunc("/remove_from_basket/{id}", RemoveFromBasket)
	r.HandleFunc("/clear_basket", ClearBasket)

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
	r.HandleFunc("/add_features", addFeatures)
	r.HandleFunc("/delete_feature", delFeature)
	r.HandleFunc("/delete_all_features", delAllFeatures)

	r.HandleFunc("/edit/categories", categories)
	r.HandleFunc("/edit/categories/{id}", showCategory)
	r.HandleFunc("/update_category", updateCategory)
	r.HandleFunc("/add_category", addCategory)
	r.HandleFunc("/delete_category", delCategory)
	r.HandleFunc("/delete_all_categories", delAllCategories)

	r.HandleFunc("/edit/units", units)
	r.HandleFunc("/edit/units/{id}", editUnit)
	r.HandleFunc("/update_unit", updateUnit)
	r.HandleFunc("/add_unit", addUnit)
	r.HandleFunc("/delete_unit", delUnit)
	r.HandleFunc("/delete_all_units", delAllUnits)

	r.HandleFunc("/add_processor", addProcessor)
	r.HandleFunc("/add_motherboard", addMotherboard)
	r.HandleFunc("/add_videocard", addVideocard)
	r.HandleFunc("/add_ram", addRAM)

	r.HandleFunc("/edit/orders", orders)

	r.HandleFunc("/categories", showCategories)

	r.HandleFunc("/categories/processors", showProcessors)

	r.HandleFunc("/categories/motherboards", showMotherboards)

	r.HandleFunc("/categories/videocards", showVideocards)

	r.HandleFunc("/categories/rams", showRams)

	r.HandleFunc("/units/{id}", showUnit)

	r.HandleFunc("/constructor", showConstructor)
	r.HandleFunc("/add_constructor/{id}", addConstructor)
	r.HandleFunc("/remove_constructor/{id}", removeConstructor)
	r.HandleFunc("/constructor/clear", clearConstructor)
	r.HandleFunc("/constructor/order", orderConstructor)

	port := configFile.HTTP.Port

	filesFolder = configFile.Files.Folder

	http.ListenAndServe(":"+port, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	processors := makeData(false, "", "Процессор", "processors", "add_basket", "Добавить в корзину", r)
	motherboards := makeData(false, "", "Материнская плата", "motherboards", "add_basket", "Добавить в корзину", r)
	videocards := makeData(false, "", "Видеокарта", "videocards", "add_basket", "Добавить в корзину", r)
	rams := makeData(false, "", "Оперативная память", "rams", "add_basket", "Добавить в корзину", r)
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "header.html", "Список товаров")
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
		if strings.Contains(page, "edit") || strings.Contains(page, "add") {
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

func checkAccess(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "cookie-name")
	isLogged, _ := session.Values["authenticated"].(bool)
	email, _ := session.Values["login"].(string)
	user, _ := database.GetUserByEmail(email)
	if !(user.AccessLevel == 10 && isLogged) {
		return false
	}
	return true
}

func hidePage(w http.ResponseWriter, r *http.Request) bool {
	access := checkAccess(w, r)
	if access == false {
		logger.Info(r.URL, "У тебя здесь нет власти!")
		http.NotFound(w, r)
		return false
	}
	return true
}

//---------------------------------------------------------------------------//

// DataFull хранит данные об товаре и категории для показа
type DataFull struct {
	ShowCategory  bool
	CategoryNames string
	CategoryLink  string
	Data          []Data
}

// Data хранит данные об товаре для показа
type Data struct {
	Picture  string
	LinkUnit string
	Name     string
	Price    int
	Link     string
	Text     string
}

func makeData(showCategory bool, categoryNames string, categoryName string, categoryLink string, actionLink string, actionText string, r *http.Request) DataFull {
	categories, _ := database.GetAllCategories()
	categoryID := 0
	for i := 0; i < len(categories); i++ {
		if categories[i].Name == categoryName {
			categoryID = categories[i].ID
		}
	}
	units, err := database.GetUnitsByCategoryID(categoryID)
	if err != nil {
		logger.Warn(err, "Не удалось загрузить список товаров категории ", categoryName, "!")
		return DataFull{}
	}

	var filteredUnits []database.Unit
	if showCategory {
		filteredUnits = filterUnits(units, categoryName, r)
	} else {
		filteredUnits = units
	}

	data := make([]Data, len(filteredUnits))
	for i := 0; i < len(filteredUnits); i++ {
		if len(filteredUnits[i].Pictures) > 0 {
			data[i].Picture = filteredUnits[i].Pictures[0]
		}
		data[i].LinkUnit = "/units/" + strconv.Itoa(filteredUnits[i].ID)
		data[i].Name = filteredUnits[i].Name
		data[i].Price = filteredUnits[i].Price
		data[i].Link = "/" + actionLink + "/" + strconv.Itoa(filteredUnits[i].ID)
		data[i].Text = actionText
	}

	dataFull := DataFull{}
	dataFull.ShowCategory = showCategory
	dataFull.CategoryNames = categoryNames
	dataFull.CategoryLink = categoryLink
	dataFull.Data = data

	return dataFull
}

func execute(w http.ResponseWriter, templateName string, data interface{}) {
	err := tpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func makeSingleData(ID int, categoryLink string, actionLink string, actionText string) Data {
	unit, err := database.GetUnit(ID)
	if err != nil {
		return Data{}
	}

	data := Data{}
	data.Picture = unit.Pictures[0]
	data.LinkUnit = "/categories/" + categoryLink + "/" + strconv.Itoa(unit.ID)
	data.Name = unit.Name
	data.Price = unit.Price
	data.Link = "/" + actionLink + "/" + strconv.Itoa(unit.ID)
	data.Text = actionText
	return data
}

// NothingData хранит данные, если товар не указан
type NothingData struct {
	Picture string
	Text    string
}

func makeNothingData(picture string, text string) NothingData {
	data := NothingData{}
	data.Picture = picture
	data.Text = text
	return data
}

func filterUnits(units []database.Unit, categoryName string, r *http.Request) []database.Unit {
	session, _ := store.Get(r, "cookie-name")
	processorID, _ := session.Values["processor"].(int)
	processor, _ := database.GetUnit(processorID)

	motherboardID, _ := session.Values["motherboard"].(int)
	motherboard, _ := database.GetUnit(motherboardID)

	videocardID, _ := session.Values["videocard"].(int)
	videocard, _ := database.GetUnit(videocardID)

	ramID, _ := session.Values["ram"].(int)
	ram, _ := database.GetUnit(ramID)
	filteredUnits := []database.Unit{}
	if len(units) != 0 {
		if categoryName == "Процессор" {
			for i := 0; i < len(units); i++ {
				if motherboard.Name != "" {
					// Сокет
					if units[i].Features[4].Value != motherboard.Features[3].Value {
						continue
					}
				}
				filteredUnits = append(filteredUnits, units[i])
			}
		}
		if categoryName == "Материнская плата" {
			for i := 0; i < len(units); i++ {
				if processor.Name != "" {
					// Сокет
					if units[i].Features[3].Value != processor.Features[4].Value {
						continue
					}
				}
				if ram.Name != "" {
					// DDR
					if units[i].Features[7].Value < ram.Features[4].Value {
						continue
					}
				}
				if videocard.Name != "" {
					// PCI-Express
					stringArray := strings.Split(videocard.Features[8].Value, " ")
					intefaceValue := stringArray[2]
					if units[i].Features[9].Value < intefaceValue {
						continue
					}
				}
				filteredUnits = append(filteredUnits, units[i])
			}
		}
		if categoryName == "Видеокарта" {
			for i := 0; i < len(units); i++ {
				if motherboard.Name != "" {
					// PCI-Express
					stringArray := strings.Split(units[i].Features[8].Value, " ")
					intefaceValue := stringArray[2]
					if intefaceValue > motherboard.Features[9].Value {
						continue
					}
				}
				filteredUnits = append(filteredUnits, units[i])
			}
		}
		if categoryName == "Оперативная память" {
			for i := 0; i < len(units); i++ {
				if motherboard.Name != "" {
					// DDR
					if units[i].Features[4].Value > motherboard.Features[7].Value {
						continue
					}
				}
				filteredUnits = append(filteredUnits, units[i])
			}
		}
	}
	return filteredUnits
}
