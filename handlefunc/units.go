package handlefunc

import (
	"net/http"
)

func showProcessors(w http.ResponseWriter, r *http.Request) {
	data := makeData("Процессор", "processors", "add_busket", "Добавить в корзину")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_units.html", data)
	}
}

func showMotherboards(w http.ResponseWriter, r *http.Request) {
	data := makeData("Материнская плата", "motherboards", "add_busket", "Добавить в корзину")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_units.html", data)
	}
}

func showVideocards(w http.ResponseWriter, r *http.Request) {
	data := makeData("Видеокарта", "videocards", "add_busket", "Добавить в корзину")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_units.html", data)
	}
}

func showRams(w http.ResponseWriter, r *http.Request) {
	data := makeData("Оперативная память", "rams", "add_busket", "Добавить в корзину")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_units.html", data)
	}
}
