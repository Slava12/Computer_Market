package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/database"
)

func showProcessors(w http.ResponseWriter, r *http.Request) {
	categories, _ := database.GetAllCategories()
	categoryID := 0
	for i := 0; i < len(categories); i++ {
		if categories[i].Name == "Процессор" {
			categoryID = categories[i].ID
		}
	}
	processors, _ := database.GetUnitsByCategoryID(categoryID)
	type Data struct {
		Picture string
		Name    string
		Price   int
		Link    string
		Text    string
	}
	//data := []Data{}
	data := make([]Data, len(processors))
	for i := 0; i < len(processors); i++ {
		data[i].Picture = processors[i].Pictures[0]
		data[i].Name = processors[i].Name
		data[i].Price = processors[i].Price
		data[i].Link = "/add_busket"
		data[i].Text = "Добавить в корзину"
	}
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "show_units.html", data)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
