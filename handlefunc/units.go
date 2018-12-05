package handlefunc

import (
	"net/http"
	"strconv"

	"github.com/Slava12/Computer_Market/database"
	"github.com/Slava12/Computer_Market/logger"
	"github.com/gorilla/mux"
)

func showProcessors(w http.ResponseWriter, r *http.Request) {
	data := makeData(false, "", "Процессор", "processors", "add_basket", "Добавить в корзину", r)
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "header.html", "Процессоры")
		execute(w, "show_units.html", data)
	}
}

func showMotherboards(w http.ResponseWriter, r *http.Request) {
	data := makeData(false, "", "Материнская плата", "motherboards", "add_basket", "Добавить в корзину", r)
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "header.html", "Материнские платы")
		execute(w, "show_units.html", data)
	}
}

func showVideocards(w http.ResponseWriter, r *http.Request) {
	data := makeData(false, "", "Видеокарта", "videocards", "add_basket", "Добавить в корзину", r)
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "header.html", "Видеокарты")
		execute(w, "show_units.html", data)
	}
}

func showRams(w http.ResponseWriter, r *http.Request) {
	data := makeData(false, "", "Оперативная память", "rams", "add_basket", "Добавить в корзину", r)
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "header.html", "Оперативная память")
		execute(w, "show_units.html", data)
	}
}

func showUnit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	unitIDstring := vars["id"]
	unitID, errString := strconv.Atoi(unitIDstring)
	if errString != nil {
		logger.Warn(errString, "Не удалось конвертировать строку в число!")
		return
	}
	unit, err := database.GetUnit(unitID)
	if err != nil {
		logger.Warn(err, "Не удалось получить данные о товаре ", unitID, "!")
		return
	}
	pairs, err := database.GetPairsByUnitID(unit.ID)
	if err != nil {
		logger.Warn(err, "Не удалось получить список пар товара ", unit.ID, "!")
		return
	}
	unitsID := make([]int, 0)
	for _, pair := range pairs {
		if pair.One != unit.ID {
			unitsID = append(unitsID, pair.One)
		} else {
			unitsID = append(unitsID, pair.Two)
		}
	}
	units := []database.Unit{}
	for _, id := range unitsID {
		tempUnit, err := database.GetUnit(id)
		if err != nil {
			logger.Warn(err, "Не удалось получить данные о товаре ", id, "!")
			return
		}
		units = append(units, tempUnit)
	}
	actionLink := "add_basket"
	actionText := "Добавить в корзину"
	data := make([]Data, len(units))
	for i := 0; i < len(units); i++ {
		if len(units[i].Pictures) > 0 {
			data[i].Picture = units[i].Pictures[0]
		}
		data[i].LinkUnit = "/units/" + strconv.Itoa(units[i].ID)
		data[i].Name = units[i].Name
		data[i].Price = units[i].Price
		data[i].Link = "/" + actionLink + "/" + strconv.Itoa(units[i].ID)
		data[i].Text = actionText
	}

	dataFull := DataFull{}
	dataFull.ShowCategory = false
	dataFull.CategoryNames = ""
	dataFull.CategoryLink = ""
	dataFull.Data = data
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "header.html", unit.Name)
		execute(w, "show_unit.html", unit)
		if len(units) > 0 {
			execute(w, "header.html", "С этим товаром покупают")
			execute(w, "show_units.html", dataFull)
		}
	}
}
