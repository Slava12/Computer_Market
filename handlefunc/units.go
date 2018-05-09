package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/errortemplate"
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

func showProcessor(w http.ResponseWriter, r *http.Request) {
	message := "Приносим свои извинения, работа над страницей ещё не завершена."
	errorMessage := errortemplate.Error{Message: message, Link: "/index"}
	execute(w, "error.html", errorMessage)
}

func showMotherboard(w http.ResponseWriter, r *http.Request) {
	message := "Приносим свои извинения, работа над страницей ещё не завершена."
	errorMessage := errortemplate.Error{Message: message, Link: "/index"}
	execute(w, "error.html", errorMessage)
}

func showVideocard(w http.ResponseWriter, r *http.Request) {
	message := "Приносим свои извинения, работа над страницей ещё не завершена."
	errorMessage := errortemplate.Error{Message: message, Link: "/index"}
	execute(w, "error.html", errorMessage)
}

func showRam(w http.ResponseWriter, r *http.Request) {
	message := "Приносим свои извинения, работа над страницей ещё не завершена."
	errorMessage := errortemplate.Error{Message: message, Link: "/index"}
	execute(w, "error.html", errorMessage)
}
