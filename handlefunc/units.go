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
	if r.Method == "GET" {
		menu(w, r)
		err := tpl.ExecuteTemplate(w, "unit.html", unit)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
