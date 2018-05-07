package handlefunc

import (
	"net/http"

	"github.com/Slava12/Computer_Market/errortemplate"
)

func showProcessors(w http.ResponseWriter, r *http.Request) {
	data := makeData("Процессор", "processors", "add_basket", "Добавить в корзину")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_units.html", data)
	}
}

func showMotherboards(w http.ResponseWriter, r *http.Request) {
	data := makeData("Материнская плата", "motherboards", "add_basket", "Добавить в корзину")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_units.html", data)
	}
}

func showVideocards(w http.ResponseWriter, r *http.Request) {
	data := makeData("Видеокарта", "videocards", "add_basket", "Добавить в корзину")
	if r.Method == "GET" {
		menu(w, r)
		execute(w, "show_units.html", data)
	}
}

func showRams(w http.ResponseWriter, r *http.Request) {
	data := makeData("Оперативная память", "rams", "add_basket", "Добавить в корзину")
	if r.Method == "GET" {
		menu(w, r)
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
